
package api

type LogParams []interface{}

func (p *LogParams) Add(field string, value interface{}) {
    *p = append(*p, field, value)
}

// Syntactic sugar for append(p, other...)
// Returns a new instance of LogParams with all values from both p and other
func (p LogParams) Combine(other LogParams) LogParams {
    return append(p, other...)
}
