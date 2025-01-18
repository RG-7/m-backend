# Managenment Backend API's

- Create User (POST)
`http://localhost:8081/auth/register`
```
{
  "email": "john.doe@example.com",
  "password": "hashedPassword123", 
  "role": "faculty",
  "mobileno": "1234567890",
  "employeeId": "EMP12345",
  "facultyCode": "F12345",
  "department": "Computer Science",
  "departmentCode": "CS01",
  "designation": "Professor",
  "availability": "M-F 9AM-12PM"
}
```

- Login User (POST)
`http://localhost:8081/auth/login`
```
{
  "email": "john.doe@example.com",
  "password": "hashedPassword123"
}
```
```
Reply :- 

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI2NzhjMzJhNzJmYzc4ODAxMTQyMDJmZWUiLCJleHAiOjE3MzkwNDEyODgsImlhdCI6MTczNzI0MTI4OH0.0qk2yd_S4AHjFYbma6FDSEigfNBq7B-Dsh1uxdwI_9o"
}
```

- Verify token (GET)
`http://localhost:8081/auth/validate`
Header Format
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI2NzhjMzJhNzJmYzc4ODAxMTQyMDJmZWUiLCJleHAiOjE3MzkwNDEyODgsImlhdCI6MTczNzI0MTI4OH0.0qk2yd_S4AHjFYbma6FDSEigfNBq7B-Dsh1uxdwI_9o

```
Reply
```
{
  "msg": "Token is valid",
  "user": {
    "id": "678c32a72fc7880114202fee",
    "email": "john.doe@example.com",
    "password": "$2a$10$39SQPKga97dh1uceIcxXieFk96Tdz0s13FYKFmI09DA88B9GDQ0.m",
    "role": "faculty",
    "mobileno": "1234567890",
    "employeeId": "EMP12345",
    "facultyCode": "F12345",
    "department": "Computer Science",
    "departmentCode": "CS01",
    "designation": "Professor",
    "availability": "M-F 9AM-12PM"
  }
}
```