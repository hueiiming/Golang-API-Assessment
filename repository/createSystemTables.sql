CREATE TABLE IF NOT EXISTS teachers (
  teacher_id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS students (
  student_id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS registrations (
   registration_id SERIAL PRIMARY KEY,
   teacher_email
   student_email
   UNIQUE (teacher_email, student_email)
);

CREATE TABLE IF NOT EXISTS suspensions (
     suspension_id SERIAL PRIMARY KEY,
     student_id INT REFERENCES students(student_id),
     UNIQUE (student_id)
);

CREATE TABLE IF NOT EXISTS notifications (
   notification_id SERIAL PRIMARY KEY,
   teacher_id INT REFERENCES teachers(teacher_id),
   notification_text TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS mentioned_students (
    mention_id SERIAL PRIMARY KEY,
    notification_id INT REFERENCES notifications(notification_id),
    student_id INT REFERENCES students(student_id)
);