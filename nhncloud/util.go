package nhncloud

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	gophercloud "github.com/nhn/nhncloud.gophercloud"
)

// BuildRequest takes an opts struct and builds a request body for
// Gophercloud to execute.
func BuildRequest(opts interface{}, parent string) (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	b = AddValueSpecs(b)

	return map[string]interface{}{parent: b}, nil
}

// CheckDeleted checks the error to see if it's a 404 (Not Found) and, if so,
// sets the resource ID to the empty string instead of throwing an error.
func CheckDeleted(d *schema.ResourceData, err error, msg string) error {
	if _, ok := err.(gophercloud.ErrDefault404); ok {
		d.SetId("")
		return nil
	}

	return fmt.Errorf("%s %s: %s", msg, d.Id(), err)
}

// GetRegion returns the region that was specified in the resource. If a
// region was not set, the provider-level region is checked. The provider-level
// region can either be set by the region argument or by OS_REGION_NAME.
func GetRegion(d *schema.ResourceData, config *Config) string {
	if v, ok := d.GetOk("region"); ok {
		return v.(string)
	}

	return config.Region
}

// AddValueSpecs expands the 'value_specs' object and removes 'value_specs'
// from the reqeust body.
func AddValueSpecs(body map[string]interface{}) map[string]interface{} {
	if body["value_specs"] != nil {
		for k, v := range body["value_specs"].(map[string]interface{}) {
			body[k] = v
		}
		delete(body, "value_specs")
	}

	return body
}

// MapValueSpecs converts ResourceData into a map.
func MapValueSpecs(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("value_specs").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}

func checkForRetryableError(err error) *resource.RetryError {
	switch e := err.(type) {
	case gophercloud.ErrDefault500:
		return resource.RetryableError(err)
	case gophercloud.ErrDefault409:
		return resource.RetryableError(err)
	case gophercloud.ErrDefault503:
		return resource.RetryableError(err)
	case gophercloud.ErrUnexpectedResponseCode:
		if e.GetStatusCode() == 504 || e.GetStatusCode() == 502 {
			return resource.RetryableError(err)
		} else {
			return resource.NonRetryableError(err)
		}
	default:
		return resource.NonRetryableError(err)
	}
}

func suppressEquivalentTimeDiffs(k, old, new string, d *schema.ResourceData) bool {
	oldTime, err := time.Parse(time.RFC3339, old)
	if err != nil {
		return false
	}

	newTime, err := time.Parse(time.RFC3339, new)
	if err != nil {
		return false
	}

	return oldTime.Equal(newTime)
}

func resourceNetworkingAvailabilityZoneHintsV2(d *schema.ResourceData) []string {
	rawAZH := d.Get("availability_zone_hints").([]interface{})
	azh := make([]string, len(rawAZH))
	for i, raw := range rawAZH {
		azh[i] = raw.(string)
	}
	return azh
}

func expandVendorOptions(vendOptsRaw []interface{}) map[string]interface{} {
	vendorOptions := make(map[string]interface{})

	for _, option := range vendOptsRaw {
		for optKey, optValue := range option.(map[string]interface{}) {
			vendorOptions[optKey] = optValue
		}
	}

	return vendorOptions
}

func expandObjectReadTags(d *schema.ResourceData, tags []string) {
	d.Set("all_tags", tags)

	allTags := d.Get("all_tags").(*schema.Set)
	desiredTags := d.Get("tags").(*schema.Set)
	actualTags := allTags.Intersection(desiredTags)
	if !actualTags.Equal(desiredTags) {
		d.Set("tags", expandToStringSlice(actualTags.List()))
	}
}

func expandObjectUpdateTags(d *schema.ResourceData) []string {
	allTags := d.Get("all_tags").(*schema.Set)
	oldTagsRaw, newTagsRaw := d.GetChange("tags")
	oldTags, newTags := oldTagsRaw.(*schema.Set), newTagsRaw.(*schema.Set)

	allTagsWithoutOld := allTags.Difference(oldTags)

	return expandToStringSlice(allTagsWithoutOld.Union(newTags).List())
}

func expandObjectTags(d *schema.ResourceData) []string {
	rawTags := d.Get("tags").(*schema.Set).List()
	tags := make([]string, len(rawTags))

	for i, raw := range rawTags {
		tags[i] = raw.(string)
	}

	return tags
}

func expandToMapStringString(v map[string]interface{}) map[string]string {
	m := make(map[string]string)
	for key, val := range v {
		if strVal, ok := val.(string); ok {
			m[key] = strVal
		}
	}

	return m
}

func expandToStringSlice(v []interface{}) []string {
	s := make([]string, len(v))
	for i, val := range v {
		if strVal, ok := val.(string); ok {
			s[i] = strVal
		}
	}

	return s
}

// strSliceContains checks if a given string is contained in a slice
// When anybody asks why Go needs generics, here you go.
func strSliceContains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func sliceUnion(a, b []string) []string {
	var res []string
	for _, i := range a {
		if !strSliceContains(res, i) {
			res = append(res, i)
		}
	}
	for _, k := range b {
		if !strSliceContains(res, k) {
			res = append(res, k)
		}
	}
	return res
}

// compatibleMicroversion will determine if an obtained microversion is
// compatible with a given microversion.
func compatibleMicroversion(direction, required, given string) (bool, error) {
	if direction != "min" && direction != "max" {
		return false, fmt.Errorf("Invalid microversion direction %s. Must be min or max", direction)
	}

	if required == "" || given == "" {
		return false, nil
	}

	requiredParts := strings.Split(required, ".")
	if len(requiredParts) != 2 {
		return false, fmt.Errorf("Not a valid microversion: %s", required)
	}

	givenParts := strings.Split(given, ".")
	if len(givenParts) != 2 {
		return false, fmt.Errorf("Not a valid microversion: %s", given)
	}

	requiredMajor, requiredMinor := requiredParts[0], requiredParts[1]
	givenMajor, givenMinor := givenParts[0], givenParts[1]

	requiredMajorInt, err := strconv.Atoi(requiredMajor)
	if err != nil {
		return false, fmt.Errorf("Unable to parse microversion: %s", required)
	}

	requiredMinorInt, err := strconv.Atoi(requiredMinor)
	if err != nil {
		return false, fmt.Errorf("Unable to parse microversion: %s", required)
	}

	givenMajorInt, err := strconv.Atoi(givenMajor)
	if err != nil {
		return false, fmt.Errorf("Unable to parse microversion: %s", given)
	}

	givenMinorInt, err := strconv.Atoi(givenMinor)
	if err != nil {
		return false, fmt.Errorf("Unable to parse microversion: %s", given)
	}

	switch direction {
	case "min":
		if requiredMajorInt == givenMajorInt {
			if requiredMinorInt <= givenMinorInt {
				return true, nil
			}
		}
	case "max":
		if requiredMajorInt == givenMajorInt {
			if requiredMinorInt >= givenMinorInt {
				return true, nil
			}
		}
	}

	return false, nil
}

func validateJSONObject(v interface{}, k string) ([]string, []error) {
	if v == nil || v.(string) == "" {
		return nil, []error{fmt.Errorf("%q value must not be empty", k)}
	}

	var j map[string]interface{}
	s := v.(string)

	err := json.Unmarshal([]byte(s), &j)
	if err != nil {
		return nil, []error{fmt.Errorf("%q must be a JSON object: %s", k, err)}
	}

	return nil, nil
}

func diffSuppressJSONObject(k, old, new string, d *schema.ResourceData) bool {
	if strSliceContains([]string{"{}", ""}, old) &&
		strSliceContains([]string{"{}", ""}, new) {
		return true
	}
	return false
}

// Metadata in openstack are not fully replaced with a "set"
// operation, instead, it's only additive, and the existing
// metadata are only removed when set to `null` value in json.
func mapDiffWithNilValues(oldMap, newMap map[string]interface{}) (output map[string]interface{}) {
	output = make(map[string]interface{})

	for k, v := range newMap {
		output[k] = v
	}

	for key := range oldMap {
		_, ok := newMap[key]
		if !ok {
			output[key] = nil
		}
	}

	return
}

// @tc-iaas-compute/1452
// Helpful facility to dispatch a structuralized data set into a variable of the other type
func NewTransformer() *Transformer {
	return &Transformer{
		regexp.MustCompile("(.)([A-Z][a-z]+)"),
		regexp.MustCompile("([a-z0-9])([A-Z])"),
	}
}

type Transformer struct {
	matchFirstCap *regexp.Regexp
	matchAllCap   *regexp.Regexp
}

func (this Transformer) isSequence(T reflect.Kind) bool {
	switch T {
	case reflect.Slice, reflect.Array:
		return true
	default:
		return false
	}
}
func (this Transformer) isPrimitive(T reflect.Kind) bool {
	switch T {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		return true
	default:
		return false
	}
}

func (this Transformer) toSnakeCase(str string) string {
	snake := this.matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = this.matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func (this Transformer) extract(v reflect.Value) interface{} {
	if this.isSequence(v.Kind()) {
		return this.transformFromArrayImplement(v)
	}
	if v.Kind() == reflect.Struct {
		return this.transformFromStructImplement(v)
	}
	return v.Interface()
}

func (this Transformer) transformFromArrayImplement(v reflect.Value) []interface{} {
	n := v.Len()
	r := make([]interface{}, n, n)
	for i := 0; i < n; i++ {
		r[i] = this.extract(v.Index(i))
	}
	return r
}

func (this Transformer) transformFromStructImplement(v reflect.Value) map[string]interface{} {
	r := make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Struct {
			this.flatten(r, v.Field(i), this.toSnakeCase(v.Type().Field(i).Name))
		} else {
			r[this.toSnakeCase(v.Type().Field(i).Name)] = this.extract(v.Field(i))
		}
	}
	return r
}

func (this Transformer) Transform(arbitrary interface{}) map[string]interface{} {
	return this.transformFromStructImplement(reflect.ValueOf(arbitrary).Elem())
}

func (this Transformer) flatten(r map[string]interface{}, v reflect.Value, s string) {
	for i := 0; i < v.NumField(); i++ {
		r[s+"_"+this.toSnakeCase(v.Type().Field(i).Name)] = v.Field(i).Interface()
	}
}
