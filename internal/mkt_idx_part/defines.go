/************************************************************
Copyright (C), shanghai, Quant360 Tech. Co., Ltd.
FileName: defines.go
Author: hfouyang
Version :1.0
Date:02012023
Description: definite some defines for project
History: hfouyang : 0201 : v_1.0 : init first
		 hfouyang : 0310 : v_0.1.5 : add feat api
***********************************************************/

// package src
package mkt_idx_part

import "unsafe"

const (
	MDSPXY_DALIY_VERSION string = "v0.1.6"
)

// const defines
const (
	MAX_WR_BUFFER_SIZE             uint32 = 4096
	TEST_CONST_DEFINES             uint32 = 0
	MDS_MAX_INSTR_CODE_LEN         uint32 = 9
	MDS_MAX_TRADING_PHASE_CODE_LEN uint32 = 9
	MDS_MAX_L2_DISCLOSE_ORDERS_CNT uint32 = 50
)

// err_code
const (
	FUNC_FAIL       int8   = -1
	FUNC_SUCCESS    int8   = 0
	BEYOND_MAX_SIZE uint32 = 0x8081
)

const (
	MDSAPI_CFG_DEFAULT_NAME           string = "mds_proxy_cfg.conf"
	MDSAPI_CFG_DEFAULT_SECTION_LOGGER string = "log"
	MDSAPI_CFG_DEFAULT_SECTION        string = "mds_client"
)

type CB func(*Tlv[any])

// MdsServerConfig ini.conf映射结构体
type MdsServerConfig struct {
	Name   string `ini:"name"`
	Passwd string `ini:"passwd"`
	Server string `ini:"server"`
}

type DataRootPath struct {
	Path string `ini:"path"`
}
type HFMDataPath struct {
	Path string `ini:"hfm_path"`
}
type DataOutputPath struct {
	Path string `ini:"data_output_path"`
}

type ProxyDataPath struct {
	Path string `ini:"proxy_data_path"`
}

// Header 消息头
type Header struct {
	MsgId   uint32
	BodyLen uint32
}

// Login 登录消息
type Login struct {
	SenderCompID     [20]byte // 发送方代码
	TargetCompID     [20]byte // 接收方代码
	HeartBtInt       int32    // 心跳间隔,单位为秒
	Password         [16]byte // 密码
	DefaultApplVerID [32]byte // 二进制协议版本,即通信版本号
}

// Tail 消息尾
type Tail struct {
	CheckSum uint32
}

//type MdsLogon struct {
//
//}

type MdsMktDataRequestReq struct {
	SubMode                 uint8 // 订阅模式
	TickType                uint8 // 数据模式
	SseStockFlag            int8  // 上证股票(股票/债券/基金)产品的订阅标志
	SseIndexFlag            int8  // 上证指数产品的订阅标志
	SseOptFlag              int8  // 上证期权产品的订阅标志
	SzseStockFlag           int8  // 深圳股票(股票/债券/基金)产品的订阅标志
	SzseIndexFlag           int8  // 深圳指数产品的订阅标志
	SzseOptFlag             int8  // 深圳期权产品的订阅标志
	IsRequireInitialMktData uint8 // 在推送实时行情数据之前, 是否需要推送已订阅产品的初始的行情快照
	ChannelNos              uint8 // 待订阅的内部频道号 (供内部使用, 尚未对外开放)
	tickExpireType          uint8 // 已废弃字段
	TickRebuildFlag         uint8 // 逐笔数据的数据重建标识 (标识是否订阅重建到的逐笔数据)
	DataTypes               int32 // 订阅的数据种类
	BeginTime               int32 // 请求订阅的行情数据的起始时间 (格式为: HHMMSS 或 HHMMSSsss)
	SubSecurityCnt          int32 // 本次订阅的产品数量 (订阅列表中的产品数量)
}

// MdsMktDataRequestEntry 行情订阅请求的订阅产品条目
type MdsMktDataRequestEntry struct {
	ExchangeId    uint8    // 交易所代码
	MdProductType uint8    // 证券类型
	filler        [2]uint8 // 内存补齐
	InstrId       int32    // 证券代码
}

// SMsgHeader 通用消息头 ==>对应api中的SMsgHeadT
type SMsgHeader struct {
	MsgFlag      uint8 //消息标志
	MsgId        uint8 // 消息代码
	Status       uint8 // 状态码
	DetailStatus uint8 // 明细状态码
	MsgSize      int32 // 消息长度
}

const (
	BUF_FILE_MAX_SIZE_10M int64 = 10240000
	BUF_FILE_MAX_SIZE_15M int64 = 15360000
)

type MdsPriceLevelEntry struct {
	Price          int32 // 委托价格 (价格单位精确到元后四位, 即: 1元=10000)
	NumberOfOrders int32 // 价位总委托笔数 (Level1不揭示该值, 固定为0)
	OrderQty       int64 // 委托数量 (@note 上交所债券的数量单位为手)
}

const (
	SH_L2_SNAPSHOT_BODY_SIZE = uint32(unsafe.Sizeof(MdsSHL2Snapshot{}))
)

type STimespec32T struct {
	Tv_sec  int32 // seconds
	Tv_nsec int32 // and nanoseconds
}

type MdsMktDataSnapshotHead struct {
	ExchId        uint8 /**< 交易所代码(沪/深) @see eMdsExchangeIdT */
	MdProductType uint8 /**< 行情产品类型 (股票/指数/期权) @see eMdsMdProductTypeT */
	IsRepeated    int8  /**< 是否是重复的行情 (内部使用, 小于0表示数据倒流) */
	OrigMdSource  uint8 /**< 原始行情数据来源 @see eMdsMsgSourceT */
	TradeDate     int32 /**< 交易日期 (YYYYMMDD, 8位整型数值) */
	UpdateTime    int32 /**< 行情时间 (HHMMSSsss, 交易所时间) */
	InstrId       int32 /**< 证券代码 (转换为整数类型的证券代码) */
	BodyLength    int16 /**< 实际数据长度 */
	BodyType      uint8 /**< 快照数据对应的消息类型 @see eMdsMsgTypeT */
	//uint8           mdStreamType;           /**< 快照数据对应的消息类型 @deprecated 已废弃, 为了兼容旧版本而暂时保留 */
	SubStreamType  uint8  /**< 行情数据类别 @see eMdsSubStreamTypeT */
	ChannelNo      uint16 /**< 频道代码 (仅适用于深交所, 对于上交所快照该字段无意义, 取值范围[0..9999]) */
	DataVersion    uint16 /**< 行情数据的更新版本号 */
	OrigTickSeq    uint32 /**< 对应的原始行情的序列号 (供内部使用) */
	DirectSourceId uint32 /**< 内部数据来源标识 (仅内部使用) */
	/** 消息原始接收时间 (从网络接收到数据的最初时间) */
	OrigNetTime   STimespec32T
	RecvTime      STimespec32T
	CollectedTime STimespec32T
	ProcessedTime STimespec32T
	PushingTime   STimespec32T
}

type MdsSHL2Snapshot struct {
	SecurityID        [MDS_MAX_INSTR_CODE_LEN]byte         // 证券代码
	TradingPhaseCode  [MDS_MAX_TRADING_PHASE_CODE_LEN]byte // 产品实时阶段及标志
	Filler            [6]byte                              // 字节对齐
	NumTrades         uint64                               // 成交笔数
	TotalVolumeTraded uint64                               /**< 成交总量 (@note 上交所债券的数量单位为手) */
	TotalValueTraded  int64                                /**< 成交总金额 (金额单位精确到元后四位, 即: 1元=10000) */
	PrevClosePx       int32                                /**< 昨日收盘价 (价格单位精确到元后四位, 即: 1元=10000) */
	OpenPx            int32                                /**< 今开盘价 (价格单位精确到元后四位, 即: 1元=10000) */
	HighPx            int32                                /**< 最高价 (价格单位精确到元后四位, 即: 1元=10000) */
	LowPx             int32                                /**< 最低价 (价格单位精确到元后四位, 即: 1元=10000) */
	TradePx           int32                                /**< 成交价 (最新价. 单位精确到元后四位, 即: 1元=10000) */
	ClosePx           int32                                /**< 今收盘价/期权收盘价 (加权平均收盘价, 适用于上交所行情和深交所债券现券交易产品. 单位精确到元后四位, 即: 1元=10000) */
	IOPV              int32                                /**< 基金份额参考净值/ETF申赎的单位参考净值 (适用于基金. 单位精确到元后四位, 即: 1元=10000) */
	NAV               int32                                /**< 基金 T-1 日净值 (适用于基金, 上交所Level-2实时行情里面没有该字段. 单位精确到元后四位, 即: 1元=10000) */
	TotalLongPosition uint64                               /**< 期权合约总持仓量 (适用于期权. 单位为张) */
	//BondWeightedAvgPx       int32                                /** 债券加权平均价 (适用于质押式回购及债券现券交易产品, 表示质押式回购成交量加权平均利率及债券现券交易成交量加权平均价. 价格单位精确到元后四位, 即: 1元=10000) */
	//BondAuctionTradePx      int32                                // 深交所债券匹配成交的最近成交价 (仅适用于深交所债券现券交易产品. 价格单位精确到元后四位, 即: 1元=10000)
	//BondAuctionVolumeTraded uint64                               /** 深交所债券匹配成交的成交总量 (仅适用于深交所债券现券交易产品) */
	TotalBidQty        int64 /**< 委托买入总量 (@note 上交所债券的数量单位为手) */
	TotalOfferQty      int64 /**< 委托卖出总量 (@note 上交所债券的数量单位为手) */
	WeightedAvgBidPx   int32 /**< 加权平均委买价格 (价格单位精确到元后四位, 即: 1元=10000) */
	WeightedAvgOfferPx int32 /**< 加权平均委卖价格 (价格单位精确到元后四位, 即: 1元=10000) */
	BidPriceLevel      int32 /**< 买方委托价位数 (实际的委托价位总数, @note 仅适用于上交所) */
	OfferPriceLevel    int32 /**< 卖方委托价位数 (实际的委托价位总数, @note 仅适用于上交所) */
	//BondAuctionValueTraded  int64                                /**< 深交所债券匹配成交的成交总金额 (@note 仅适用于深交所债券现券交易产品. 金额单位精确到元后四位, 即: 1元=10000) */
	BidLevels   [10]MdsPriceLevelEntry // 十档买盘价位信息
	OfferLevels [10]MdsPriceLevelEntry // 十档卖盘价位信息
}

type MdsMktL2Snapshot struct {
	Header   MdsMktDataSnapshotHead
	Snapshot MdsSHL2Snapshot
}

type MdsMktSZL2Snapshot struct {
	Header   MdsMktDataSnapshotHead
	Snapshot MdsSZL2Snapshot
}

type MdsSZL2Snapshot struct {
	SecurityID             [MDS_MAX_INSTR_CODE_LEN]byte         // 证券代码
	TradingPhaseCode       [MDS_MAX_TRADING_PHASE_CODE_LEN]byte // 产品实时阶段及标志
	Filler                 [6]byte                              // 字节对齐
	NumTrades              uint64                               // 成交笔数
	TotalVolumeTraded      uint64                               /**< 成交总量 (@note 上交所债券的数量单位为手) */
	TotalValueTraded       int64                                /**< 成交总金额 (金额单位精确到元后四位, 即: 1元=10000) */
	PrevClosePx            int32                                /**< 昨日收盘价 (价格单位精确到元后四位, 即: 1元=10000) */
	OpenPx                 int32                                /**< 今开盘价 (价格单位精确到元后四位, 即: 1元=10000) */
	HighPx                 int32                                /**< 最高价 (价格单位精确到元后四位, 即: 1元=10000) */
	LowPx                  int32                                /**< 最低价 (价格单位精确到元后四位, 即: 1元=10000) */
	TradePx                int32                                /**< 成交价 (最新价. 单位精确到元后四位, 即: 1元=10000) */
	ClosePx                int32                                /**< 今收盘价/期权收盘价 (加权平均收盘价, 适用于上交所行情和深交所债券现券交易产品. 单位精确到元后四位, 即: 1元=10000) */
	IOPV                   int32                                /**< 基金份额参考净值/ETF申赎的单位参考净值 (适用于基金. 单位精确到元后四位, 即: 1元=10000) */
	NAV                    int32                                /**< 基金 T-1 日净值 (适用于基金, 上交所Level-2实时行情里面没有该字段. 单位精确到元后四位, 即: 1元=10000) */
	TotalLongPosition      uint64                               /**< 期权合约总持仓量 (适用于期权. 单位为张) */
	TotalBidQty            int64                                /**< 委托买入总量 (@note 上交所债券的数量单位为手) */
	TotalOfferQty          int64                                /**< 委托卖出总量 (@note 上交所债券的数量单位为手) */
	WeightedAvgBidPx       int32                                /**< 加权平均委买价格 (价格单位精确到元后四位, 即: 1元=10000) */
	WeightedAvgOfferPx     int32                                /**< 加权平均委卖价格 (价格单位精确到元后四位, 即: 1元=10000) */
	BondAuctionValueTraded int64                                /**< 深交所债券匹配成交的成交总金额 (@note 仅适用于深交所债券现券交易产品. 金额单位精确到元后四位, 即: 1元=10000) */
	BidLevels              [10]MdsPriceLevelEntry               // 十档买盘价位信息
	OfferLevels            [10]MdsPriceLevelEntry               // 十档卖盘价位信息
}
type MdsL2BestOrdersSnapshot struct {
	/** 证券代码 C6 / C8 (如: '600000' 等) */
	SecurityID        [MDS_MAX_INSTR_CODE_LEN]byte
	Filler            [5]uint8 /**< 按64位对齐的填充域 */
	NoBidOrders       uint8    /**< 买一价位的揭示委托笔数 */
	NoOfferOrders     uint8    /**< 卖一价位的揭示委托笔数 */
	TotalVolumeTraded uint64   /**< 成交总量 (来自快照行情的冗余字段) */
	BestBidPrice      int32    /**< 最优申买价 (价格单位精确到元后四位, 即: 1元=10000) */
	BestOfferPrice    int32    /**< 最优申卖价 (价格单位精确到元后四位, 即: 1元=10000) */
	/** 买一价位的委托明细(前50笔) */
	BidOrderQty [MDS_MAX_L2_DISCLOSE_ORDERS_CNT]int32

	/** 卖一价位的委托明细(前50笔) */
	OfferOrderQty [MDS_MAX_L2_DISCLOSE_ORDERS_CNT]int32
}

type MdsL1Snapshot struct {
	SecurityID        [MDS_MAX_INSTR_CODE_LEN]byte         // 证券代码
	TradingPhaseCode  [MDS_MAX_TRADING_PHASE_CODE_LEN]byte // 产品实时阶段及标志
	Filler            [6]byte                              // 字节对齐
	NumTrades         uint64                               // 成交笔数
	TotalVolumeTraded uint64                               /**< 成交总量 (@note 上交所债券的数量单位为手) */
	TotalValueTraded  int64                                /**< 成交总金额 (金额单位精确到元后四位, 即: 1元=10000) */
	PrevClosePx       int32                                /**< 昨日收盘价 (价格单位精确到元后四位, 即: 1元=10000) */
	OpenPx            int32                                /**< 今开盘价 (价格单位精确到元后四位, 即: 1元=10000) */
	HighPx            int32                                /**< 最高价 (价格单位精确到元后四位, 即: 1元=10000) */
	LowPx             int32                                /**< 最低价 (价格单位精确到元后四位, 即: 1元=10000) */
	TradePx           int32                                /**< 成交价 (最新价. 单位精确到元后四位, 即: 1元=10000) */
	ClosePx           int32                                /**< 今收盘价/期权收盘价 (加权平均收盘价, 适用于上交所行情和深交所债券现券交易产品. 单位精确到元后四位, 即: 1元=10000) */
	IOPV              int32                                /**< 基金份额参考净值/ETF申赎的单位参考净值 (适用于基金. 单位精确到元后四位, 即: 1元=10000) */
	NAV               int32                                /**< 基金 T-1 日净值 (适用于基金, 上交所Level-2实时行情里面没有该字段. 单位精确到元后四位, 即: 1元=10000) */
	TotalLongPosition uint64                               /**< 期权合约总持仓量 (适用于期权. 单位为张) */
	//BondWeightedAvgPx       int32                        /** 债券加权平均价 (适用于质押式回购及债券现券交易产品, 表示质押式回购成交量加权平均利率及债券现券交易成交量加权平均价. 价格单位精确到元后四位, 即: 1元=10000) */
	//BondAuctionTradePx      int32                        // 深交所债券匹配成交的最近成交价 (仅适用于深交所债券现券交易产品. 价格单位精确到元后四位, 即: 1元=10000)
	//BondAuctionVolumeTraded uint64                       /** 深交所债券匹配成交的成交总量 (仅适用于深交所债券现券交易产品) */
	TotalBidQty        int64 /**< 委托买入总量 (@note 上交所债券的数量单位为手) */
	TotalOfferQty      int64 /**< 委托卖出总量 (@note 上交所债券的数量单位为手) */
	WeightedAvgBidPx   int32 /**< 加权平均委买价格 (价格单位精确到元后四位, 即: 1元=10000) */
	WeightedAvgOfferPx int32 /**< 加权平均委卖价格 (价格单位精确到元后四位, 即: 1元=10000) */
	BidPriceLevel      int32 /**< 买方委托价位数 (实际的委托价位总数, @note 仅适用于上交所) */
	OfferPriceLevel    int32 /**< 卖方委托价位数 (实际的委托价位总数, @note 仅适用于上交所) */
	//BondAuctionValueTraded  int64                        /**< 深交所债券匹配成交的成交总金额 (@note 仅适用于深交所债券现券交易产品. 金额单位精确到元后四位, 即: 1元=10000) */
	BidLevels   [5]MdsPriceLevelEntry // 五档买盘价位信息
	OfferLevels [5]MdsPriceLevelEntry // 五档卖盘价位信息
}

type IndexSnapshotBody struct {
	SecurityID        [MDS_MAX_INSTR_CODE_LEN]byte         // 证券代码
	TradingPhaseCode  [MDS_MAX_TRADING_PHASE_CODE_LEN]byte // 产品实时阶段及标志
	Filler            [6]byte                              // 字节对齐
	NumTrades         uint64                               // 成交笔数
	TotalVolumeTraded uint64                               /**< 成交总量 (@note 上交所债券的数量单位为手) */
	TotalValueTraded  int64                                /**< 成交总金额 (金额单位精确到元后四位, 即: 1元=10000) */
	PrevCloseIdx      int32                                /**< 昨日收盘价 (价格单位精确到元后四位, 即: 1元=10000) */
	OpenIdx           int32                                /**< 今开盘价 (价格单位精确到元后四位, 即: 1元=10000) */
	HighIdx           int32                                /**< 最高价 (价格单位精确到元后四位, 即: 1元=10000) */
	LowIdx            int32                                /**< 最低价 (价格单位精确到元后四位, 即: 1元=10000) */
	TradeIdx          int32                                /**< 成交价 (最新价. 单位精确到元后四位, 即: 1元=10000) */
	CloseIdx          int32                                /**< 今收盘价/期权收盘价 (加权平均收盘价, 适用于上交所行情和深交所债券现券交易产品. 单位精确到元后四位, 即: 1元=10000) */
	StockNum          int32                                // 统计量指标样本个数 (@note 仅适用于深交所成交量统计指标)
	Filler1           int32                                // 字节对齐
}

type MdsL2Trade struct {
	ExchId          uint8  /**< 交易所代码(沪/深) @see eMdsExchangeIdT */
	MdProductType   uint8  /**< 行情产品类型 (股票) @see eMdsMdProductTypeT */
	IsRepeated      int8   /**< 是否是重复的行情 (内部使用, 小于0表示回补的逐笔重建数据) */
	OrigMdSource    uint8  /**< 原始行情数据来源 @see eMdsMsgSourceT */
	TradeDate       int32  /**< 交易日期 (YYYYMMDD, 非官方数据) */
	TransactTime    int32  /**< 成交时间 (HHMMSSsss) */
	InstrId         int32  /**< 证券代码 (转换为整数类型的证券代码) */
	ChannelNo       uint16 /**< 频道代码 [0..9999] */
	Reserve         uint16 /**< 按64位对齐的保留字段 */
	ApplSeqNum      uint32
	SecurityID      [MDS_MAX_INSTR_CODE_LEN]byte // 证券代码
	ExecType        byte
	TradeBSFlag     byte
	SubStreamType   uint8
	SseBizIndex     uint32
	Filler          uint64
	TradePrice      int32
	TradeQty        int32
	TradeMoney      int64
	BidApplSeqNum   int64
	OfferApplSeqNum int64
	OrigNetTime     STimespec32T
	RecvTime        STimespec32T
	CollectedTime   STimespec32T
	ProcessedTime   STimespec32T
	PushingTime     STimespec32T
}

type MdsL2Order struct {
	ExchId        uint8  /**< 交易所代码(沪/深) @see eMdsExchangeIdT */
	MdProductType uint8  /**< 行情产品类型 (股票) @see eMdsMdProductTypeT */
	IsRepeated    int8   /**< 是否是重复的行情 (内部使用, 小于0表示回补的逐笔重建数据) */
	OrigMdSource  uint8  /**< 原始行情数据来源 @see eMdsMsgSourceT */
	TradeDate     int32  /**< 交易日期 (YYYYMMDD, 非官方数据) */
	TransactTime  int32  /**< 成交时间 (HHMMSSsss) */
	InstrId       int32  /**< 证券代码 (转换为整数类型的证券代码) */
	ChannelNo     uint16 /**< 频道代码 [0..9999] */
	Reserve       uint16 /**< 按64位对齐的保留字段 */
	ApplSeqNum    uint32
	SecurityID    [MDS_MAX_INSTR_CODE_LEN]byte // 证券代码
	Side          byte
	OrderType     byte
	SubStreamType uint8
	SseBizIndex   uint32
	SseOrderNo    int64
	Price         int32
	OrderQty      int32
	OrigNetTime   STimespec32T
	RecvTime      STimespec32T
	CollectedTime STimespec32T
	ProcessedTime STimespec32T
	PushingTime   STimespec32T
}

type AllInOne struct {
	Ticks     []MdsL2Trade
	Orders    []MdsL2Order
	Snapshots []MdsMktL2Snapshot
}

type SZAllInOne struct {
	ticks     []MdsL2Trade
	orders    []MdsL2Order
	snapshots []MdsMktSZL2Snapshot
}

const (
	ESecurityID int = iota
	EDateTime
	EPreClosePx
	EOpenPx
	EHighPx
	ELowPx
	ELastPx
	ETotalVolumeTrade
	ETotalValueTrade
	EInstrumentStatus
	EBidPrice
	EBidOrderQty
	EBidNumOrders
	EBidOrders
	EOfferPrice
	EOfferOrderQty
	EOfferNumOrders
	EOfferOrders
	ENumTrades
	EIOPV
	ETotalBidQty
	ETotalOfferQty
	EWeightedAvgBidPx
	EWeightedAvgOfferPx
	ETotalBidNumber
	ETotalOfferNumber
	EBidTradeMaxDuration
	EOfferTradeMaxDuration
	ENumBidOrders
	ENumOfferOrders
)

const (
	TICK_TYPE uint32 = iota + 1
	ORDER_TYPE
	SNAPSHOT_TYPE
	SZTICK_TYPE
	SZORDER_TYPE
	SZSNAPSHOT_TYPE
	ALLINONE_TYPE
)

type CacheInter interface {
	MdsL2Order | MdsL2Trade | MdsMktSZL2Snapshot | any
}

type Tlv[T CacheInter] struct {
	Header TlvHeader
	Data   T
}

type TlvHeader struct {
	DataTyp uint32
	BufLen  uint32
}

const (
	EOSecurityID int = iota
	EOTransactTime
	EOOrderNo
	EOPrice
	EOBalance
	EOOrderBSFlag
	EOOrdType
	EOOrderIndex
	EOChannelNo
	EOBizIndex
)

const (
	ETSecurityID int = iota
	ETTradeTime
	ETTradePrice
	ETTradeQty
	ETTradeAmount
	ETBuyNo
	ETSellNo
	ETTradeIndex
	ETChannelNo
	ETTradeBSFlag
	ETBizIndex
)

type TestTlv struct {
	Typ uint16
	Siz uint16
	Dat []byte
}

// SZ Trade enum
const (
	ESZApplSeqNum int = iota
	ESZBidApplSeqNum
	ESZSendingTime
	ESZPrice
	ESZChannelNo
	ESZQty
	ESZOfferApplSeqNum
	ESZAmt
	ESZExecType
	ESZTransactTime
)

// SZ Order enum
const (
	ESZORDOrderQty int = iota
	ESZORDOrdType
	ESZORDTransactTime
	ESZORDExpirationDays
	ESZORDSide
	ESZORDApplSeqNum
	ESZORDContactor
	ESZORDSendingTime
	ESZORDPrice
	ESZORDChannelNo
	ESZORDExpirationType
	ESZORDContactInfo
	ESZORDConfirmID
)

// SZ Snapshot enum
const (
	SNAPESZSendingTime int = iota
	ESZMsgSeqNum
	ESZImageStatus
	ESZQuotTime
	ESZPreClosePx
	ESZOpenPx
	ESZHighPx
	ESZLowPx
	ESZLastPx
	ESZClosePx
	ESZVolume
	ESZAmount
	ESZAveragePx
	ESZBidPrice
	ESZBidOrderQty
	ESZBidNumOrders
	ESZBidOrders
	ESZOfferPrice
	ESZOfferOrderQty
	ESZOfferNumOrders
	ESZOfferOrders
	ESZNumTrades
	ESZTotalBidQty
	ESZWeightedAvgBidPx
	ESZTotalOfferQty
	ESZWeightedAvgOfferPx
	ESZChange1
	ESZChange2
	ESZTotalLongPosition
	ESZPeRatio1
	ESZPeRatio2
	ESZUpperLimitPx
	ESZLowerLimitPx
	ESZWeightedAvgPxChg
	ESZPreWeightedAvgPx
	ESZTradingPhaseCode
	ESZNoOrdersB1
	ESZNoOrdersS1
)

const (
	SH_FLAG int = iota
	SZ_FLAG
)
