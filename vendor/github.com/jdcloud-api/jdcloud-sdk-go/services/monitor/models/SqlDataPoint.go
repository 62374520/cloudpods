// Copyright 2018 JDCLOUD.COM
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// NOTE: This class is auto generated by the jdcloud code generator program.

package models


type SqlDataPoint struct {

    /* 目前统一用jcloud  */
    AppCode string `json:"appCode"`

    /* 资源的类型，取值sqlserver  */
    ServiceCode string `json:"serviceCode"`

    /* 资源所在的地域  */
    Region string `json:"region"`

    /* 资源的uuid  */
    ResourceId string `json:"resourceId"`

    /* 监控指标名称，长度不超过255字节，只允许英文、数字、下划线_、点.,  [0-9][a-z] [A-Z] [. _ ]， 其它会返回err  */
    Metric string `json:"metric"`

    /* 毫秒级时间戳，早于当前时间30天的不能写入；建议的上报时间戳：上报时间间隔的整数倍，如上报间隔为5ms，则建议上报的时间戳为 time = current timestamp - (current timestamp % time interval) = 1487647187007 - （1487647187007 % 5） = 1487647187007 -2 = 1487647187005  */
    Time int64 `json:"time"`

    /* 上报的监控值，即慢sql语句已经执行的时间(单位s)  */
    Value int64 `json:"value"`

    /* SQL开始执行的时间  */
    Start_time string `json:"start_time"`

    /* SQL已执行时间(单位s)  */
    Execution_time int64 `json:"execution_time"`

    /* 会话ID  */
    Session_id string `json:"session_id"`

    /* 数据库库名  */
    Database string `json:"database"`

    /* 客户端IP地址  */
    Client_net_address string `json:"client_net_address"`

    /* 用户名  */
    Loginname string `json:"loginname"`

    /* SQL会话请求状态  */
    Status string `json:"status"`

    /* SQL详细文本  */
    Sqlstr string `json:"sqlstr"`
}