package orm

import "fmt"

func ExampleCondition() {
	dialect := &TestDialect{}
	cond := NewCondition(dialect).
		Appendf("name = ?", "Medivh").
		AppendMap(map[string]interface{}{"age": 20, "male": true}).
		AppendIn("pet", "cat", "dog", "tiger")
	t := New(dialect).
		Select("user", "id", "name", "age", "pet").
		WhereTemplate(cond.And().Bracket()).
		Where("foo = ?", "bar").
		Build()
	fmt.Println(t)
	// Output: 
	// "select 'id','name','age','pet' from 'user' where (name = ? and age = ? and male = ? and 'pet' in (?, ?, ?)) and foo = ?": []interface {}{"Medivh", 20, true, "cat", "dog", "tiger", "bar", "Medivh", 20, true, "cat", "dog", "tiger", "bar"}
}
