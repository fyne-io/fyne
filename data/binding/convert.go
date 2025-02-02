package binding

import (
	"fyne.io/fyne/v2"
)

// BoolToString creates a binding that connects a Bool data item to a String.
// Changes to the Bool will be pushed to the String and setting the string will parse and set the
// Bool if the parse was successful.
//
// Since: 2.0
func BoolToString(v Bool) String {
	return toStringComparable[bool](v, formatBool, parseBool)
}

// BoolToStringWithFormat creates a binding that connects a Bool data item to a String and is
// presented using the specified format. Changes to the Bool will be pushed to the String and setting
// the string will parse and set the Bool if the string matches the format and its parse was successful.
//
// Since: 2.0
func BoolToStringWithFormat(v Bool, format string) String {
	return toStringWithFormatComparable[bool](v, format, "%t", formatBool, parseBool)
}

// FloatToString creates a binding that connects a Float data item to a String.
// Changes to the Float will be pushed to the String and setting the string will parse and set the
// Float if the parse was successful.
//
// Since: 2.0
func FloatToString(v Float) String {
	return toStringComparable[float64](v, formatFloat, parseFloat)
}

// FloatToStringWithFormat creates a binding that connects a Float data item to a String and is
// presented using the specified format. Changes to the Float will be pushed to the String and setting
// the string will parse and set the Float if the string matches the format and its parse was successful.
//
// Since: 2.0
func FloatToStringWithFormat(v Float, format string) String {
	return toStringWithFormatComparable[float64](v, format, "%f", formatFloat, parseFloat)
}

// IntToFloat creates a binding that connects an Int data item to a Float.
//
// Since: 2.5
func IntToFloat(val Int) Float {
	v := &fromIntTo[float64]{from: val, parser: internalFloatToInt, formatter: internalIntToFloat}
	val.AddListener(v)
	return v
}

// FloatToInt creates a binding that connects a Float data item to an Int.
//
// Since: 2.5
func FloatToInt(v Float) Int {
	i := &toInt[float64]{from: v, parser: internalFloatToInt, formatter: internalIntToFloat}
	v.AddListener(i)
	return i
}

// IntToString creates a binding that connects a Int data item to a String.
// Changes to the Int will be pushed to the String and setting the string will parse and set the
// Int if the parse was successful.
//
// Since: 2.0
func IntToString(v Int) String {
	return toStringComparable[int](v, formatInt, parseInt)
}

// IntToStringWithFormat creates a binding that connects a Int data item to a String and is
// presented using the specified format. Changes to the Int will be pushed to the String and setting
// the string will parse and set the Int if the string matches the format and its parse was successful.
//
// Since: 2.0
func IntToStringWithFormat(v Int, format string) String {
	return toStringWithFormatComparable[int](v, format, "%d", formatInt, parseInt)
}

// URIToString creates a binding that connects a URI data item to a String.
// Changes to the URI will be pushed to the String and setting the string will parse and set the
// URI if the parse was successful.
//
// Since: 2.1
func URIToString(v URI) String {
	return toString[fyne.URI](v, uriToString, compareURI, uriFromString)
}

// StringToBool creates a binding that connects a String data item to a Bool.
// Changes to the String will be parsed and pushed to the Bool if the parse was successful, and setting
// the Bool update the String binding.
//
// Since: 2.0
func StringToBool(str String) Bool {
	v := &fromStringTo[bool]{from: str, formatter: parseBool, parser: formatBool}
	str.AddListener(v)
	return v
}

// StringToBoolWithFormat creates a binding that connects a String data item to a Bool and is
// presented using the specified format. Changes to the Bool will be parsed and if the format matches and
// the parse is successful it will be pushed to the String. Setting the Bool will push a formatted value
// into the String.
//
// Since: 2.0
func StringToBoolWithFormat(str String, format string) Bool {
	if format == "%t" { // Same as not using custom format.
		return StringToBool(str)
	}

	v := &fromStringTo[bool]{from: str, format: format}
	str.AddListener(v)
	return v
}

// StringToFloat creates a binding that connects a String data item to a Float.
// Changes to the String will be parsed and pushed to the Float if the parse was successful, and setting
// the Float update the String binding.
//
// Since: 2.0
func StringToFloat(str String) Float {
	v := &fromStringTo[float64]{from: str, formatter: parseFloat, parser: formatFloat}
	str.AddListener(v)
	return v
}

// StringToFloatWithFormat creates a binding that connects a String data item to a Float and is
// presented using the specified format. Changes to the Float will be parsed and if the format matches and
// the parse is successful it will be pushed to the String. Setting the Float will push a formatted value
// into the String.
//
// Since: 2.0
func StringToFloatWithFormat(str String, format string) Float {
	if format == "%f" { // Same as not using custom format.
		return StringToFloat(str)
	}

	v := &fromStringTo[float64]{from: str, format: format}
	str.AddListener(v)
	return v
}

// StringToInt creates a binding that connects a String data item to a Int.
// Changes to the String will be parsed and pushed to the Int if the parse was successful, and setting
// the Int update the String binding.
//
// Since: 2.0
func StringToInt(str String) Int {
	v := &fromStringTo[int]{from: str, parser: formatInt, formatter: parseInt}
	str.AddListener(v)
	return v
}

// StringToIntWithFormat creates a binding that connects a String data item to a Int and is
// presented using the specified format. Changes to the Int will be parsed and if the format matches and
// the parse is successful it will be pushed to the String. Setting the Int will push a formatted value
// into the String.
//
// Since: 2.0
func StringToIntWithFormat(str String, format string) Int {
	if format == "%d" { // Same as not using custom format.
		return StringToInt(str)
	}

	v := &fromStringTo[int]{from: str, format: format}
	str.AddListener(v)
	return v
}

// StringToURI creates a binding that connects a String data item to a URI.
// Changes to the String will be parsed and pushed to the URI if the parse was successful, and setting
// the URI update the String binding.
//
// Since: 2.1
func StringToURI(str String) URI {
	v := &fromStringTo[fyne.URI]{from: str, parser: uriToString, formatter: uriFromString}
	str.AddListener(v)
	return v
}
