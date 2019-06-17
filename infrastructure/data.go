package infrastructure

type Cred struct {
	Host   string
	Port   string
	DBName string
}

var Credentials = Cred{Host:"mongodb://127.0.0.1", Port:"27017", DBName:"StarProject",}
var UsersCollectionName, GroupCollectionName, TabelCollectionName,
	SchoolsCollectionName, PricesCollectionName, ChildrenCollectionName  = "users", "groups", "tabel", "schools", "prices", "children"
