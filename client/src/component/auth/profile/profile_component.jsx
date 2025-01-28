import React, { useEffect, useState } from 'react';
import axios from 'axios';

const ProfileComponent = () => {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const token = localStorage.getItem('token');
        if (!token) {
          throw new Error('Токен отсутствует');
        }

        const response = await axios.get('http://localhost:5000/profile', {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        setUser(response.data);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchUserData();
  }, []);

  if (loading) {
    return <div>Загрузка...</div>;
  }

  if (error) {
    return <div>Ошибка: {error}</div>;
  }

  return (
    <div>
      <h1>Профиль пользователя</h1>
      {user && (
        <div>
          <p>Полное имя: {user.fullName}</p>
          <p>Имя пользователя: {user.username}</p>
          <p>Email: {user.email}</p>
          <p>Дата рождения: {user.birthDate}</p>
          <p>Место жительства: {user.residence}</p>
          <p>Образование: {user.education}</p>
          <p>Место рождения: {user.birthPlace}</p>
          <p>Домашний адрес: {user.homeAddress}</p>
          <p>Паспортные данные: {user.passportData}</p>
          <p>СНИЛС: {user.snils}</p>
        </div>
      )}
    </div>
  );
};

export default ProfileComponent;