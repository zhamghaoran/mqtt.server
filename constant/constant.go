package constant

const (
	Reversed       int = 0  //保留
	CONNECT        int = 1  // 客户端请求连接服务端
	CONNACK        int = 2  //连接确认报文确定
	PUBLISH        int = 3  //发布消息
	PUBACK         int = 4  //Qos1消息发布确认
	PUBREC         int = 5  //发布收到保证第一步
	PUBREL         int = 6  //发布释放 保证交付第二部
	PUBCOMP        int = 7  //Qos2消息发布完成  保证交付第二部
	SUBSCRIBE      int = 8  //客户端订阅发布
	SUBACK         int = 9  //订阅请求报文确定
	UNSUBSCRIBE    int = 10 //客户端取消订阅
	UNSUBSCRIBEACK int = 11 //客户端取消订阅请求
	PINGREG        int = 12 //心跳请求
	PINGRESQ       int = 13 //心跳响应
	DISCONNECT     int = 14 //客户端断开连接
	TYPELEN        int = 15
)
