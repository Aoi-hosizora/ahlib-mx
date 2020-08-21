# xstatus

### DbStatus

+ `type DbStatus int8`
+ `DbSuccess`
+ `DbNotFound`
+ `DbExisted`
+ `DbFailed`
+ `DbTagA`
+ `DbTagB`
+ `DbTagC`
+ `DbTagD`
+ `DbTagE`
+ `(d DbStatus) String() string`

### FsmStatus

+ `type FsmStatus int8`
+ `FsmNone`
+ `FsmInState`
+ `FsmFinal`
+ `(f FsmStatus) String() string`

### JwtStatus

+ `type JwtStatus int8`
+ `JwtExpired`
+ `JwtIssuer`
+ `JwtNotIssued`
+ `JwtNotValid`
+ `JwtID`
+ `JwtAudience`
+ `JwtSubject`
+ `JwtInvalid`
+ `JwtUserErr`
+ `JwtFailed`
+ `JwtTagA`
+ `JwtTagB`
+ `JwtTagC`
+ `JwtTagD`
+ `JwtTagE`
+ `(j JwtStatus) String() string`
