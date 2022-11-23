package utils

//  JsonGetInt64
//    performs float64 (json numbers are always float64)
//    type assertion and casts to int64.
//
//    if the type assertion fails, the function defaults 0 (zero).
// ---------------------------------------------------------

func JsonGetInt64(input interface{}) int64 {
	v, res := input.(float64)
	if res {
		return (int64)(v)
	}
	return 0
}
