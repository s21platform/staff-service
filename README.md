# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [api/staff.proto](#api_staff-proto)
    - [ChangePasswordIn](#staff-ChangePasswordIn)
    - [ChangePasswordOut](#staff-ChangePasswordOut)
    - [CheckAuthIn](#staff-CheckAuthIn)
    - [CheckAuthOut](#staff-CheckAuthOut)
    - [CreateIn](#staff-CreateIn)
    - [CreateOut](#staff-CreateOut)
    - [DeleteIn](#staff-DeleteIn)
    - [DeleteOut](#staff-DeleteOut)
    - [GetIn](#staff-GetIn)
    - [GetOut](#staff-GetOut)
    - [ListIn](#staff-ListIn)
    - [ListOut](#staff-ListOut)
    - [LoginIn](#staff-LoginIn)
    - [LoginOut](#staff-LoginOut)
    - [LogoutIn](#staff-LogoutIn)
    - [LogoutOut](#staff-LogoutOut)
    - [Permissions](#staff-Permissions)
    - [RefreshTokenIn](#staff-RefreshTokenIn)
    - [RefreshTokenOut](#staff-RefreshTokenOut)
    - [Staff](#staff-Staff)
    - [UpdateIn](#staff-UpdateIn)
    - [UpdateOut](#staff-UpdateOut)
  
    - [StaffService](#staff-StaffService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="api_staff-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## api/staff.proto



<a name="staff-ChangePasswordIn"></a>

### ChangePasswordIn
Запрос на смену пароля


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| old_password | [string](#string) |  |  |
| new_password | [string](#string) |  |  |
| access_token | [string](#string) |  | access_token используется для идентификации пользователя |






<a name="staff-ChangePasswordOut"></a>

### ChangePasswordOut
Ответ на смену пароля


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  |  |






<a name="staff-CheckAuthIn"></a>

### CheckAuthIn
Запрос на проверку авторизации


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_token | [string](#string) |  |  |






<a name="staff-CheckAuthOut"></a>

### CheckAuthOut
Ответ на проверку авторизации


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| authorized | [bool](#bool) |  |  |
| staff | [Staff](#staff-Staff) |  |  |






<a name="staff-CreateIn"></a>

### CreateIn
Запрос на создание сотрудника


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| login | [string](#string) |  |  |
| password | [string](#string) |  |  |
| role_id | [int32](#int32) |  |  |
| permissions | [Permissions](#staff-Permissions) |  |  |






<a name="staff-CreateOut"></a>

### CreateOut
Ответ с информацией о созданном сотруднике


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| staff | [Staff](#staff-Staff) |  |  |






<a name="staff-DeleteIn"></a>

### DeleteIn
Запрос на удаление сотрудника


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="staff-DeleteOut"></a>

### DeleteOut
Ответ на удаление сотрудника


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  |  |






<a name="staff-GetIn"></a>

### GetIn
Запрос на получение информации о сотруднике


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |






<a name="staff-GetOut"></a>

### GetOut
Ответ с информацией о сотруднике


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| staff | [Staff](#staff-Staff) |  |  |






<a name="staff-ListIn"></a>

### ListIn
Запрос на получение списка сотрудников


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| page | [int32](#int32) |  |  |
| page_size | [int32](#int32) |  |  |
| search_term | [string](#string) | optional | поиск по логину |
| role_id | [int32](#int32) | optional | фильтр по роли |






<a name="staff-ListOut"></a>

### ListOut
Ответ со списком сотрудников


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| staff | [Staff](#staff-Staff) | repeated |  |
| total_count | [int32](#int32) |  |  |
| page_count | [int32](#int32) |  |  |






<a name="staff-LoginIn"></a>

### LoginIn
Запрос на авторизацию


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| login | [string](#string) |  |  |
| password | [string](#string) |  |  |






<a name="staff-LoginOut"></a>

### LoginOut
Ответ на успешную авторизацию


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_token | [string](#string) |  |  |
| refresh_token | [string](#string) |  |  |
| expires_at | [int64](#int64) |  | время истечения токена в unix timestamp |
| staff | [Staff](#staff-Staff) |  |  |






<a name="staff-LogoutIn"></a>

### LogoutIn
Запрос на выход из системы


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_token | [string](#string) |  |  |






<a name="staff-LogoutOut"></a>

### LogoutOut
Ответ на выход из системы


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  |  |






<a name="staff-Permissions"></a>

### Permissions
Структура разрешений сотрудника


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access | [string](#string) | repeated |  |






<a name="staff-RefreshTokenIn"></a>

### RefreshTokenIn
Запрос на обновление токена


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| refresh_token | [string](#string) |  |  |






<a name="staff-RefreshTokenOut"></a>

### RefreshTokenOut
Ответ с обновленным токеном


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| access_token | [string](#string) |  |  |
| refresh_token | [string](#string) |  |  |
| expires_at | [int64](#int64) |  | время истечения токена в unix timestamp |






<a name="staff-Staff"></a>

### Staff
Структура данных сотрудника


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| login | [string](#string) |  |  |
| role_id | [int32](#int32) |  |  |
| role_name | [string](#string) |  |  |
| permissions | [Permissions](#staff-Permissions) |  |  |
| created_at | [int64](#int64) |  |  |
| updated_at | [int64](#int64) |  |  |






<a name="staff-UpdateIn"></a>

### UpdateIn
Запрос на обновление информации о сотруднике


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| login | [string](#string) | optional |  |
| role_id | [int32](#int32) | optional |  |
| permissions | [Permissions](#staff-Permissions) | optional |  |






<a name="staff-UpdateOut"></a>

### UpdateOut
Ответ с обновленной информацией о сотруднике


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| staff | [Staff](#staff-Staff) |  |  |





 

 

 


<a name="staff-StaffService"></a>

### StaffService
Комбинированный сервис для управления персоналом и авторизацией

=== Методы управления персоналом ===

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Get | [GetIn](#staff-GetIn) | [GetOut](#staff-GetOut) | Получение информации о сотруднике по ID |
| Create | [CreateIn](#staff-CreateIn) | [CreateOut](#staff-CreateOut) | Создание нового сотрудника |
| Update | [UpdateIn](#staff-UpdateIn) | [UpdateOut](#staff-UpdateOut) | Обновление информации о сотруднике |
| Delete | [DeleteIn](#staff-DeleteIn) | [DeleteOut](#staff-DeleteOut) | Удаление сотрудника |
| List | [ListIn](#staff-ListIn) | [ListOut](#staff-ListOut) | Получение списка сотрудников с фильтрацией и пагинацией |
| Login | [LoginIn](#staff-LoginIn) | [LoginOut](#staff-LoginOut) | Авторизация сотрудника по логину и паролю |
| RefreshToken | [RefreshTokenIn](#staff-RefreshTokenIn) | [RefreshTokenOut](#staff-RefreshTokenOut) | Обновление токена сессии |
| Logout | [LogoutIn](#staff-LogoutIn) | [LogoutOut](#staff-LogoutOut) | Выход из системы и завершение сессии |
| CheckAuth | [CheckAuthIn](#staff-CheckAuthIn) | [CheckAuthOut](#staff-CheckAuthOut) | Проверка текущего статуса авторизации |
| ChangePassword | [ChangePasswordIn](#staff-ChangePasswordIn) | [ChangePasswordOut](#staff-ChangePasswordOut) | Изменение пароля авторизованного пользователя |

 



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

