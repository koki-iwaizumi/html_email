# html_email
MySQLからユーザーのデータを取得し、そのユーザーのメールへ10秒ごとにhtmlメールを送信

## Dependency
go version go1.8

## Setup
GoとMySQLのSetupは省略
```shell
$ go get github.com/go-xorm/xorm
```

## Usage
1. MySQLに下記を登録します。
```sql
CREATE TABLE IF NOT EXISTS `person` (
	`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
	`company` varchar(255) NOT NULL DEFAULT '',
	`email` varchar(255) NOT NULL DEFAULT '',
	`name` varchar(255) NOT NULL DEFAULT '',
	`honorific` varchar(255) NOT NULL DEFAULT '',
	`post_h` varchar(3) NOT NULL DEFAULT '',
	`post_l` varchar(4) NOT NULL DEFAULT '',
	`prefecture` varchar(255) NOT NULL DEFAULT '',
	`address_h` varchar(25) NOT NULL DEFAULT '',
	`address_l` varchar(25) NOT NULL DEFAULT '',
	`jinto` varchar(10) NOT NULL DEFAULT '',
	`saibaru` varchar(10) NOT NULL DEFAULT '',
	`created` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`),
	UNIQUE KEY `company` (`company`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT INTO `person` (`id`, `company`, `email`, `name`, `honorific`, `post_h`, `post_l`, `prefecture`, `address_h`, `address_l`, `jinto`, `saibaru`, `created`) VALUES
(1, 'test_company', 'info@test.test', 'testname', '様', '111', '1111', 'prefecture', 'address_h', 'address_l', '未送信', '未送信', '2019-05-14 23:21:54');
```

2. email設定の変更
main.go
```Golang
const (
	EMAIL_HOST                = "*"
	EMAIL_PORT                = "*"
	EMAIL_USER                = "*"
	EMAIL_PASSWORD            = "*"
	EMAIL_FROMNAME            = "*"
	EMAIL_FROMADDRESS         = "*"
	EMAIL_SUBJECT             = "IoTデータロガーのご紹介！"
	EMAIL_PATH                = "view/email.html"
	EMAIL_TARGET_STATUS       = "未送信"
	EMAIL_TARGET_STATUS_AFTER = "メール送信済み"
)
```

3. db設定の変更
model/base.go
```Golang
const (
	MODEL_DRIVER   = "mysql"
	MODEL_USER     = "*"
	MODEL_PASSWORD = "*"
	MODEL_NAME     = "*"
)
```

4. 実行
```shell
$ go run main.go
```

# References
* https://blog.kannart.co.jp/coding/1093/
* https://qiita.com/cyabane/items/b0cbc9bc7526c56f5724