# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/v0/staff.v0.proto](#api_v0_staff-v0-proto)
    - [ChangePasswordRequest](#api-v0-ChangePasswordRequest)
    - [ChangePasswordResponse](#api-v0-ChangePasswordResponse)
    - [CheckAuthRequest](#api-v0-CheckAuthRequest)
    - [CheckAuthResponse](#api-v0-CheckAuthResponse)
    - [CreateStaffRequest](#api-v0-CreateStaffRequest)
    - [CreateStaffResponse](#api-v0-CreateStaffResponse)
    - [DeleteStaffRequest](#api-v0-DeleteStaffRequest)
    - [DeleteStaffResponse](#api-v0-DeleteStaffResponse)
    - [GetStaffRequest](#api-v0-GetStaffRequest)
    - [GetStaffResponse](#api-v0-GetStaffResponse)
    - [ListStaffRequest](#api-v0-ListStaffRequest)
    - [ListStaffResponse](#api-v0-ListStaffResponse)
    - [LoginRequest](#api-v0-LoginRequest)
    - [LoginResponse](#api-v0-LoginResponse)
    - [LogoutRequest](#api-v0-LogoutRequest)
    - [LogoutResponse](#api-v0-LogoutResponse)
    - [Permissions](#api-v0-Permissions)
    - [RefreshTokenRequest](#api-v0-RefreshTokenRequest)
    - [RefreshTokenResponse](#api-v0-RefreshTokenResponse)
    - [Staff](#api-v0-Staff)
    - [UpdateStaffRequest](#api-v0-UpdateStaffRequest)
    - [UpdateStaffResponse](#api-v0-UpdateStaffResponse)
  
    - [StaffService](#api-v0-StaffService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api_v0_staff-v0-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/v0/staff.v0.proto



<a name="api-v0-ChangePasswordRequest"></a>

### ChangePasswordRequest
Запрос на смену пароля


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| old_password | [string](#string) |  |  |
| new_password | [string](#string) |  |  |
| access_token | [string](#string) |  | access_token используется для идентификации пользователя |






<a name="api-v0-ChangePasswordResponse"></a>

### ChangePasswordResponse
Ответ на смену пароля


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  |  |






<a name="api-v0-CheckAuthRequest"></a>

### CheckAuthRequest
Запрос на проверку авторизации


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_token | [string](#string) |  |  |






<a name="api-v0-CheckAuthResponse"></a>

### CheckAuthResponse
Ответ на проверку авторизации


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| authorized | [bool](#bool) |  |  |
| staff | [Staff](#api-v0-Staff) |  |  |






<a name="api-v0-CreateStaffRequest"></a>

### CreateStaffRequest
Запрос на создание сотрудника


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| login | [string](#string) |  |  |
| password | [string](#string) |  |  |
| role_id | [int32](#int32) |  |  |
| permissions | [Permissions](#api-v0-Permissions) |  |  |






<a name="api-v0-CreateStaffResponse"></a>

### CreateStaffResponse
Ответ с информацией о созданном сотруднике


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| staff | [Staff](#api-v0-Staff) |  |  |






<a name="api-v0-DeleteStaffRequest"></a>

### DeleteStaffRequest
Запрос на удаление сотрудника


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="api-v0-DeleteStaffResponse"></a>

### DeleteStaffResponse
Ответ на удаление сотрудника


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  |  |






<a name="api-v0-GetStaffRequest"></a>

### GetStaffRequest
Запрос на получение информации о сотруднике


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="api-v0-GetStaffResponse"></a>

### GetStaffResponse
Ответ с информацией о сотруднике


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| staff | [Staff](#api-v0-Staff) |  |  |






<a name="api-v0-ListStaffRequest"></a>

### ListStaffRequest
Запрос на получение списка сотрудников


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| page | [int32](#int32) |  |  |
| page_size | [int32](#int32) |  |  |
| search_term | [string](#string) | optional | поиск по логину |
| role_id | [int32](#int32) | optional | фильтр по роли |






<a name="api-v0-ListStaffResponse"></a>

### ListStaffResponse
Ответ со списком сотрудников


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| staff | [Staff](#api-v0-Staff) | repeated |  |
| total_count | [int32](#int32) |  |  |
| page_count | [int32](#int32) |  |  |






<a name="api-v0-LoginRequest"></a>

### LoginRequest
Запрос на авторизацию


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| login | [string](#string) |  |  |
| password | [string](#string) |  |  |






<a name="api-v0-LoginResponse"></a>

### LoginResponse
Ответ на успешную авторизацию


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_token | [string](#string) |  |  |
| refresh_token | [string](#string) |  |  |
| expires_at | [int64](#int64) |  | время истечения токена в unix timestamp |
| staff | [Staff](#api-v0-Staff) |  |  |






<a name="api-v0-LogoutRequest"></a>

### LogoutRequest
Запрос на выход из системы


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_token | [string](#string) |  |  |






<a name="api-v0-LogoutResponse"></a>

### LogoutResponse
Ответ на выход из системы


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  |  |






<a name="api-v0-Permissions"></a>

### Permissions
Структура разрешений сотрудника


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access | [string](#string) | repeated |  |






<a name="api-v0-RefreshTokenRequest"></a>

### RefreshTokenRequest
Запрос на обновление токена


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| refresh_token | [string](#string) |  |  |






<a name="api-v0-RefreshTokenResponse"></a>

### RefreshTokenResponse
Ответ с обновленным токеном


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_token | [string](#string) |  |  |
| refresh_token | [string](#string) |  |  |
| expires_at | [int64](#int64) |  | время истечения токена в unix timestamp |






<a name="api-v0-Staff"></a>

### Staff
Структура данных сотрудника


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| login | [string](#string) |  |  |
| role_id | [int32](#int32) |  |  |
| role_name | [string](#string) |  |  |
| permissions | [Permissions](#api-v0-Permissions) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="api-v0-UpdateStaffRequest"></a>

### UpdateStaffRequest
Запрос на обновление информации о сотруднике


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| login | [string](#string) | optional |  |
| role_id | [int32](#int32) | optional |  |
| permissions | [Permissions](#api-v0-Permissions) | optional |  |






<a name="api-v0-UpdateStaffResponse"></a>

### UpdateStaffResponse
Ответ с обновленной информацией о сотруднике


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| staff | [Staff](#api-v0-Staff) |  |  |





 

 

 


<a name="api-v0-StaffService"></a>

### StaffService
Комбинированный сервис для управления персоналом и авторизацией

=== Методы управления персоналом ===

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| GetStaff | [GetStaffRequest](#api-v0-GetStaffRequest) | [GetStaffResponse](#api-v0-GetStaffResponse) | Получение информации о сотруднике по ID |
| CreateStaff | [CreateStaffRequest](#api-v0-CreateStaffRequest) | [CreateStaffResponse](#api-v0-CreateStaffResponse) | Создание нового сотрудника |
| UpdateStaff | [UpdateStaffRequest](#api-v0-UpdateStaffRequest) | [UpdateStaffResponse](#api-v0-UpdateStaffResponse) | Обновление информации о сотруднике |
| DeleteStaff | [DeleteStaffRequest](#api-v0-DeleteStaffRequest) | [DeleteStaffResponse](#api-v0-DeleteStaffResponse) | Удаление сотрудника |
| ListStaff | [ListStaffRequest](#api-v0-ListStaffRequest) | [ListStaffResponse](#api-v0-ListStaffResponse) | Получение списка сотрудников с фильтрацией и пагинацией |
| Login | [LoginRequest](#api-v0-LoginRequest) | [LoginResponse](#api-v0-LoginResponse) | Авторизация сотрудника по логину и паролю |
| RefreshToken | [RefreshTokenRequest](#api-v0-RefreshTokenRequest) | [RefreshTokenResponse](#api-v0-RefreshTokenResponse) | Обновление токена сессии |
| Logout | [LogoutRequest](#api-v0-LogoutRequest) | [LogoutResponse](#api-v0-LogoutResponse) | Выход из системы и завершение сессии |
| CheckAuth | [CheckAuthRequest](#api-v0-CheckAuthRequest) | [CheckAuthResponse](#api-v0-CheckAuthResponse) | Проверка текущего статуса авторизации |
| ChangePassword | [ChangePasswordRequest](#api-v0-ChangePasswordRequest) | [ChangePasswordResponse](#api-v0-ChangePasswordResponse) | Изменение пароля авторизованного пользователя |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

