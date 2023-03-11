module github.com/go-gourd/mysql

go 1.20

require (
	github.com/go-gourd/gourd v0.2.1
	gorm.io/gorm v1.24.6
)

//本地开发框架包替换
replace github.com/go-gourd/gourd => ../gourd

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
)
