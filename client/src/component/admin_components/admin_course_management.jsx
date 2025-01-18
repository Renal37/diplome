import React, { useState } from 'react';

const AdminCourseManagement = () => {
  const [courseName, setCourseName] = useState('');
  const [coursePrice, setCoursePrice] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    const response = await fetch('http://localhost:5000/add-course', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ name: courseName, price: coursePrice }),
    });
    if (response.ok) {
      alert('Course added successfully');
    } else {
      alert('Failed to add course');
    }
  };

  return (
    <div>
      <h1>Hello admin!</h1>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="Course Name"
          value={courseName}
          onChange={(e) => setCourseName(e.target.value)}
        />
        <input
          type="text"
          placeholder="Course Price"
          value={coursePrice}
          onChange={(e) => setCoursePrice(e.target.value)}
        />
        <button type="submit">Add Course</button>
      </form>
    </div>
  );
};

export default AdminCourseManagement;