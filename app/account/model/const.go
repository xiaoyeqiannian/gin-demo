package model

// 系统固定角色ID
const ROLE_ROOT_ID = 1              // 系统root角色,此角色只能有一个用户使用
const ROLE_ADMIN_ID = 2             // 默认管理员角色，注册用户为自己组织的管理员，

// 系统固定组织ID
const GROUP_SYS_ADMIN_ID = 1 // 总部组织ID

// 组织类型
const GROUP_KIND_PERSONAL = 1 // 个人
const GROUP_KIND_GROUP = 2    // 团体/企业

// 验证码
const VERIFY_CODE_LENGTH = 6
const VERIFY_CODE_EXPIRE_TIME = 60*5