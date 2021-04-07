package goorm

type ONString string

func ON(foreignKey string, localKey string) ONString {
	return ONString("ON " + foreignKey + "=" + localKey)
}
