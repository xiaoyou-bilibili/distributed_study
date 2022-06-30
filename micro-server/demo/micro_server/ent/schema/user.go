package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Annotations 配置注解信息
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "users"}, // 比如这里我们设置一下表的名字
	}
}

// Fields 用户表的字段信息
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique(),
		field.String("name"),
		field.Int("age"),
		field.Int("sex"),
	}
}

// Edges 实体的边缘关系，比如所属的用户组等信息
func (User) Edges() []ent.Edge {
	return nil
}

// Indexes 可以设置索引
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("id").Unique(),
	}
}
