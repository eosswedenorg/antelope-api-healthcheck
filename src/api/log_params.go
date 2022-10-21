
package api

type LogParams []interface{}

func (p *LogParams) Add(field string, value interface{}) {
    *p = append(*p, field, value)
}
