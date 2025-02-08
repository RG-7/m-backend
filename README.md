# Managenment Backend API's

- Create User (POST)
  `http://localhost:8081/auth/register`

```
{
  "email": "john.doe@example.com", 
  "password": "hashedPassword123",
  "name":"Jhon", 
  "role": "faculty",
  "mobileno": "1234567890",
  "employeeId": "EMP12345",
  "facultyCode": " ",
  "department": "Computer Science",
  "departmentCode": "CS01",
  "designation": "Professor",
  "availability": " "
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
    "name":"Jhon"
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

- GetAllUser (GET)
  `http://localhost:8081/auth/all`

- DeleteUser (GET)
  `http://localhost:8081/auth/{id}` Id in Params!

```
Example:-
http://localhost:8081/auth/678c32a72fc7880114202fee
```

- Add Lectures for a subgroup of similar range i.e. (1Q1A-1Q1E) or Single also i.e. (2INSQW) (POST) 
  `http://localhost:8081/tt/linrange`

```
{
  "courseCode": "UCS550",
  "courseName": "Network Defense",
  "facultyCode": "RPN",
  "venue": "LT-401",
  "subgroup": "1Q1A-1Q1E",
  "department": "Computer Science & Engineering",
  "startDate": "2025-01-06",
  "endDate": "2025-05-18",
  "day": "Monday",
  "time": "08:00 AM",
  "type": "l"
}
```

- Delete Lectures for a subgroup of similar range (1Q1A-1Q1E)  i.e.  or Single also i.e. (2INSQW) (POST)
  `http://localhost:8081/tt/linrangedel`

```
{
  "courseCode": "UCS550",
  "courseName": "Network Defense",
  "facultyCode": "RPN",
  "venue": "LT-401",
  "subgroup": "1Q1A-1Q1E",
  "department": "Computer Science & Engineering",
  "startDate": "2025-01-06",
  "endDate": "2025-05-18",
  "day": "Monday",
  "time": "08:00 AM",
  "type": "l"
}
```

-------------------------------- COmmon -------------------------------------
- Add Lecture/Practical/Tutorial with commonsubject and subgrous 
`http://localhost:8081/tt/c`
```
{
  "courseCode": "UCS550",
  "courseName": "Network Defense",
  "facultyCode": "RPN",
  "venue": "LT-401",
  "subgroup": ["3Q2E","3Q2D","3C2A","3M2C"],
  "department": "Computer Science & Engineering",
  "startDate": "2025-01-06",
  "endDate": "2025-05-18",
  "day": "Friday",
  "time": "11:20 AM",
  "type": "p"
}
```


- Delete Lecture/Practical/Tutorial with commonsubject and subgrous 
`http://localhost:8081/tt/cdel`
```
{
  "courseCode": "UTA024",
  "courseName": "Network Defense",
  "facultyCode": "RPN",
  "venue": "LT-401",
  "subgroup": ["3Q2E","3Q2D","3C2A","3M2C"],
  "department": "Computer Science & Engineering",
  "startDate": "2025-01-06",
  "endDate": "2025-05-18",
  "day": "Friday",
  "time": "11:20 AM",
  "type": "p"
}
```

--------------------------------- Get Timetable ----------------------------------
- Get timetable by subgroup & Date 
```http://localhost:8081/tt/subgroup/1A1A/2025-01-06```
Reply
`
[
  {
    "ID": "678d47f2c3d7b70149d739c8",
    "CourseCode": "UHU003",
    "CourseName": " ",
    "FacultyCode": "KBJ",
    "Venue": "LP-101",
    "Subgroup": "1A1A",
    "Department": " ",
    "Time": "09:40 AM",
    "Date": "2025-01-06",
    "Duration": 50,
    "Type": "L",
    "CreatedAt": "2025-01-19T18:44:02.384Z"
  },
  {
    "ID": "678d4837c3d7b70149d73a73",
    "CourseCode": "UMA023",
    "CourseName": " ",
    "FacultyCode": "ARM",
    "Venue": "LP-101",
    "Subgroup": "1A1A",
    "Department": " ",
    "Time": "10:30 AM",
    "Date": "2025-01-06",
    "Duration": 50,
    "Type": "L",
    "CreatedAt": "2025-01-19T18:45:11.732Z"
  },
  {
    "ID": "678d48b7c3d7b70149d73b44",
    "CourseCode": "UES102",
    "CourseName": " ",
    "FacultyCode": "RTK",
    "Venue": "W/SHOP",
    "Subgroup": "1A1A",
    "Department": " ",
    "Time": "01:50 PM",
    "Date": "2025-01-06",
    "Duration": 100,
    "Type": "P",
    "CreatedAt": "2025-01-19T18:47:19.147Z"
  },
  {
    "ID": "678d49edc3d7b70149d73c28",
    "CourseCode": "UES101",
    "CourseName": " ",
    "FacultyCode": "NK",
    "Venue": "F-302",
    "Subgroup": "1A1A",
    "Department": " ",
    "Time": "03:30 PM",
    "Date": "2025-01-06",
    "Duration": 50,
    "Type": "T",
    "CreatedAt": "2025-01-19T18:52:29.058Z"
  }
]
`