"webkit" {}

# Method http方法类型 #
Method -> Str

GET     :Method : "GET"
POST    :Method : "POST"
PUT     :Method : "PUT"
DELETE  :Method : "DELETE"
PATCH   :Method : "PATCH"
OPTIONS :Method : "OPTIONS"

# New_Method Method构建函数 #
New Method(v:Str) -> (r:Method) {
	<- (Method(v))
}

(me:Method) String() -> (r:Str) {
	<- (string(me))
}

# isMethod 判断是否存在的方法 #
is Method(m:Method) -> (r:Bool) {
    b := False
	m ? GET {
        b = True
	} POST {
        b = True
    } PUT {
        b = True
    } DELETE {
        b = True
    } PATCH {
        b = True
    } OPTIONS {
		b = True
	}
	<- (b)
}
