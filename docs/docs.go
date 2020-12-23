// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/trailer_api/app/get_new_app": {
            "post": {
                "description": "获取最新的 APP 版本信息",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "APP"
                ],
                "summary": "UPDATE APP",
                "operationId": "UPDATE APP",
                "parameters": [
                    {
                        "description": "UPDATE_APP",
                        "name": "name",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/app_service.AppParam"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/trailer_api/stats/record_device": {
            "post": {
                "description": "设备信息上报",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Stats"
                ],
                "summary": "设备信息上报",
                "operationId": "Record Device",
                "parameters": [
                    {
                        "description": "Device",
                        "name": "name",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/stats_service.Device"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/trailer_api/stats/record_sdk_err": {
            "post": {
                "description": "SDK 错误信息上报",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Stats"
                ],
                "summary": "SDK 错误信息上报",
                "operationId": "Record SdkError",
                "parameters": [
                    {
                        "description": "SdkError",
                        "name": "name",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/stats_service.SdkError"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/trailer_api/stats/record_sdk_event": {
            "post": {
                "description": "SDK 事件统计   参数：事件 json 数组",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Stats"
                ],
                "summary": "SDK 事件统计",
                "operationId": "Insert SdkEvent",
                "parameters": [
                    {
                        "description": "Events",
                        "name": "name",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/stats_service.SdkEvent"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/trailer_api/test/check_interface": {
            "get": {
                "description": "测试接口",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Test"
                ],
                "summary": "Test Interface",
                "operationId": "Test",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/trailer_api/trailer/get_trailer_list": {
            "post": {
                "description": "获取预告片信息",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Trailer"
                ],
                "summary": "Get TrailerList",
                "operationId": "Get TrailerList",
                "parameters": [
                    {
                        "description": "TrailerListParam",
                        "name": "TrailerListParam",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/gredis.TrailerListParam"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        },
        "/trailer_api/trailer/sync_asset": {
            "post": {
                "description": "同步 Asset 素材信息",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Trailer"
                ],
                "summary": "Sync Asset",
                "operationId": "Sync Asset",
                "parameters": [
                    {
                        "description": "SyncAsset",
                        "name": "name",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/gredis.Asset"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/app.Response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "app.Response": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer"
                },
                "data": {
                    "type": "object"
                },
                "msg": {
                    "type": "string"
                }
            }
        },
        "app_service.AppParam": {
            "type": "object",
            "required": [
                "android_version_code",
                "android_version_name",
                "channel_code",
                "cpu_arch",
                "device_model",
                "device_no",
                "sdk_name",
                "sdk_version_name",
                "storage_space"
            ],
            "properties": {
                "android_version_code": {
                    "description": "Android 版本 Code",
                    "type": "string"
                },
                "android_version_name": {
                    "description": "Android 版本名称",
                    "type": "string"
                },
                "app_name": {
                    "description": "APP 名称",
                    "type": "string"
                },
                "app_version_code": {
                    "description": "APP 版本 Code",
                    "type": "string"
                },
                "app_version_name": {
                    "description": "APP 版本名称",
                    "type": "string"
                },
                "channel_code": {
                    "description": "渠道码",
                    "type": "string"
                },
                "cpu_arch": {
                    "description": "CPU 架构",
                    "type": "string"
                },
                "device_model": {
                    "description": "设备型号",
                    "type": "string"
                },
                "device_no": {
                    "description": "设备号",
                    "type": "string"
                },
                "is_hot_update": {
                    "description": "是否热更  0:否  1:是",
                    "type": "integer"
                },
                "language": {
                    "description": "语言",
                    "type": "string"
                },
                "resolution": {
                    "description": "分辨率",
                    "type": "string"
                },
                "sdk_name": {
                    "description": "SDK 名称",
                    "type": "string"
                },
                "sdk_version_code": {
                    "description": "SDK 版本 Code",
                    "type": "string"
                },
                "sdk_version_name": {
                    "description": "SDK 版本名称",
                    "type": "string"
                },
                "storage_space": {
                    "description": "存储空间",
                    "type": "string"
                }
            }
        },
        "gredis.Asset": {
            "type": "object",
            "required": [
                "act_type",
                "duration_end_date",
                "duration_start_date",
                "id",
                "last_update_time",
                "shelf_status",
                "type"
            ],
            "properties": {
                "act_long_movie_url": {
                    "description": "长视频 url",
                    "type": "string"
                },
                "act_open_apps": {
                    "description": "需要下载打开的应用  （json）",
                    "type": "object"
                },
                "act_pop_time": {
                    "description": "二维码自动弹出时间，单位：秒",
                    "type": "integer"
                },
                "act_qrcode_bg_url": {
                    "description": "二维码背景图url",
                    "type": "string"
                },
                "act_qrcode_org_url": {
                    "description": "二维码原链接url",
                    "type": "string"
                },
                "act_qrcode_url": {
                    "description": "二维码地址",
                    "type": "string"
                },
                "act_toast": {
                    "description": "OK键引导文案",
                    "type": "string"
                },
                "act_type": {
                    "description": "OK键动作类型 1-无动作 2-打开/下载应用 3-弹出二维码 4-加载长视频",
                    "type": "integer"
                },
                "channel_code": {
                    "description": "对应的渠道 - 全部为 ALL - json数组",
                    "type": "string"
                },
                "cover_url": {
                    "description": "封面图地址",
                    "type": "string"
                },
                "del_flag": {
                    "description": "是否删除  0:否  1:是",
                    "type": "integer"
                },
                "display_order": {
                    "description": "排序",
                    "type": "integer",
                    "example": 0
                },
                "duration_end_date": {
                    "description": "资源有效期 - 结束时间",
                    "type": "string"
                },
                "duration_start_date": {
                    "description": "资源有效期 - 开始时间",
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "img_stay_time": {
                    "description": "单张图片停留时长(秒)",
                    "type": "integer"
                },
                "last_update_time": {
                    "description": "最后更新时间 - 排序使用",
                    "type": "string"
                },
                "movie_url": {
                    "description": "视频地址",
                    "type": "string"
                },
                "name": {
                    "description": "名称",
                    "type": "string"
                },
                "pic_urls": {
                    "description": "多张图片url  （json）",
                    "type": "object"
                },
                "priority": {
                    "description": "优先级  0:优先调用  1:优先下载",
                    "type": "integer"
                },
                "remark": {
                    "description": "描述",
                    "type": "string"
                },
                "score": {
                    "description": "评分",
                    "type": "object"
                },
                "screen_control_count": {
                    "description": "周期内频控次数",
                    "type": "integer",
                    "example": 0
                },
                "screen_control_cycle": {
                    "description": "频控周期  天数  自然日",
                    "type": "integer",
                    "example": 0
                },
                "shelf_status": {
                    "description": "上架状态 1-未上架 2-已上架 3-已下架",
                    "type": "integer"
                },
                "type": {
                    "description": "类型 1-视频 2-图片",
                    "type": "integer"
                },
                "view_cities": {
                    "description": "地区限制 （json）",
                    "type": "string"
                },
                "view_limit": {
                    "description": "青少年观影限制 0  - 不限制 1 – 限制",
                    "type": "integer"
                },
                "view_tags": {
                    "description": "数据标签 （json）",
                    "type": "string"
                }
            }
        },
        "gredis.TrailerListParam": {
            "type": "object",
            "required": [
                "channel_code",
                "device_no",
                "page"
            ],
            "properties": {
                "channel_code": {
                    "description": "渠道码",
                    "type": "string"
                },
                "device_no": {
                    "description": "设备号",
                    "type": "string"
                },
                "page": {
                    "description": "页码",
                    "type": "integer",
                    "example": 1
                },
                "page_size": {
                    "description": "每页数量",
                    "type": "integer",
                    "example": 20
                }
            }
        },
        "stats_service.Device": {
            "type": "object",
            "required": [
                "android_version_code",
                "android_version_name",
                "app_name",
                "app_version_code",
                "app_version_name",
                "channel_code",
                "device_model",
                "device_name",
                "device_no",
                "device_vendor",
                "ip",
                "resolution"
            ],
            "properties": {
                "android_version_code": {
                    "description": "Android 版本 Code",
                    "type": "string"
                },
                "android_version_name": {
                    "description": "Android 版本名称",
                    "type": "string"
                },
                "app_name": {
                    "description": "APP 名称",
                    "type": "string"
                },
                "app_version_code": {
                    "description": "APP 版本 Code",
                    "type": "string"
                },
                "app_version_name": {
                    "description": "APP 版本名称",
                    "type": "string"
                },
                "channel_code": {
                    "description": "渠道码",
                    "type": "string"
                },
                "device_model": {
                    "description": "设备型号",
                    "type": "string"
                },
                "device_name": {
                    "description": "设备名称",
                    "type": "string"
                },
                "device_no": {
                    "description": "设备号",
                    "type": "string"
                },
                "device_vendor": {
                    "description": "设备厂商",
                    "type": "string"
                },
                "ip": {
                    "description": "IP",
                    "type": "string"
                },
                "resolution": {
                    "description": "分辨率",
                    "type": "string"
                }
            }
        },
        "stats_service.SdkError": {
            "type": "object",
            "required": [
                "app_name",
                "app_version_code",
                "app_version_name",
                "channel_code",
                "crash_log",
                "crash_time",
                "device_no",
                "sdk_name",
                "sdk_version_code",
                "sdk_version_name"
            ],
            "properties": {
                "app_name": {
                    "description": "APP 名称",
                    "type": "string"
                },
                "app_version_code": {
                    "description": "APP 版本 Code",
                    "type": "string"
                },
                "app_version_name": {
                    "description": "APP 版本名称",
                    "type": "string"
                },
                "channel_code": {
                    "description": "渠道码",
                    "type": "string"
                },
                "crash_log": {
                    "description": "Crash 日志",
                    "type": "string"
                },
                "crash_time": {
                    "description": "Crash 时间",
                    "type": "string"
                },
                "device_no": {
                    "description": "设备号",
                    "type": "string"
                },
                "ext": {
                    "description": "自定义数据",
                    "type": "string"
                },
                "sdk_name": {
                    "description": "SDK 名称",
                    "type": "string"
                },
                "sdk_version_code": {
                    "description": "SDK 版本 Code",
                    "type": "string"
                },
                "sdk_version_name": {
                    "description": "SDK 版本名称",
                    "type": "string"
                },
                "user_id": {
                    "description": "USER ID",
                    "type": "integer"
                }
            }
        },
        "stats_service.SdkEvent": {
            "type": "object",
            "required": [
                "app_name",
                "app_version_code",
                "app_version_name",
                "channel_code",
                "device_model",
                "device_no",
                "net_type",
                "newsession_id",
                "os_version_code",
                "os_version_name",
                "screen_height",
                "screen_width",
                "sdk_name",
                "sdk_version_code",
                "sdk_version_name"
            ],
            "properties": {
                "app_name": {
                    "description": "APP 名称",
                    "type": "string"
                },
                "app_version_code": {
                    "description": "APP 版本 Code",
                    "type": "string"
                },
                "app_version_name": {
                    "description": "APP 版本名称",
                    "type": "string"
                },
                "channel_code": {
                    "description": "渠道码",
                    "type": "string"
                },
                "client_time": {
                    "description": "客户端时间 格式：2020-12-12 12:12:12",
                    "type": "string"
                },
                "device_brand": {
                    "description": "设备品牌",
                    "type": "string"
                },
                "device_model": {
                    "description": "设备型号",
                    "type": "string"
                },
                "device_no": {
                    "description": "设备号",
                    "type": "string"
                },
                "event_info": {
                    "description": "参数和参数值数据 -json数组 - 格式：[{\"event_kv_json\":{\"key1\": \"val1\", \"key2\": \"val2\"},\"event_name\":\"sss\",\"event_type\":1},...]",
                    "type": "string"
                },
                "event_kv_json": {
                    "description": "参数和参数值数据 -json数组 - 格式：{\"key1\": \"val1\", \"key2\": \"val2\"}",
                    "type": "string"
                },
                "event_name": {
                    "description": "事件名称",
                    "type": "string"
                },
                "imei": {
                    "description": "IMEI",
                    "type": "string"
                },
                "ip": {
                    "description": "IP",
                    "type": "string"
                },
                "net_type": {
                    "description": "网络类型\tWIFI/4G/5G",
                    "type": "string"
                },
                "newevent_type": {
                    "description": "事件类型  0:自定义事件,1:预置事件",
                    "type": "integer"
                },
                "newpuid": {
                    "description": "IDFA \t\t\t\t\t\tstring ` + "`" + `json:\"idfa\" ` + "`" + `",
                    "type": "string"
                },
                "newsession_id": {
                    "description": "会话 ID",
                    "type": "string"
                },
                "os_version_code": {
                    "description": "操作系统版本 Code",
                    "type": "string"
                },
                "os_version_name": {
                    "description": "操作系统版本名称",
                    "type": "string"
                },
                "page_name": {
                    "description": "所在页面",
                    "type": "string"
                },
                "screen_height": {
                    "description": "分辨率 高度",
                    "type": "integer"
                },
                "screen_width": {
                    "description": "分辨率 宽度",
                    "type": "integer"
                },
                "sdk_name": {
                    "description": "SDK 名称",
                    "type": "string"
                },
                "sdk_version_code": {
                    "description": "SDK 版本 Code",
                    "type": "string"
                },
                "sdk_version_name": {
                    "description": "SDK 版本名称",
                    "type": "string"
                },
                "signature": {
                    "description": "签名",
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
