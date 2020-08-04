# xstatus

### DbStatus

+ `type DbStatus uint8`
+ `DbSuccess`
+ `DbNotFound`
+ `DbExisted`
+ `DbFailed`
+ `(d DbStatus) String() string`

### FsmStatus

+ `type FsmStatus uint8`
+ `FsmNone`
+ `FsmInState`
+ `FsmFinal`
+ `(f FsmStatus) String() string`
