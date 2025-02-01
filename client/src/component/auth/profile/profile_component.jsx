import { useState, useEffect } from "react";
import "./profile_component.css";


const Profile = () => {
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);


  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const response = await fetch("http://localhost:5000/profile", {
          method: "GET",
          credentials: "include",
        });

        if (!response.ok) {
          throw new Error("Ошибка при загрузке профиля");
        }

        const data = await response.json();
        setProfile(data);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, []);

  if (loading) {
    return <p>Загрузка...</p>;
  }

  if (error) {
    return <p>{error}</p>;
  }

  return profile ? (
    <div className="profile-container">
      <h1 className="profile-header">Ваш профиль {profile.username}</h1>
      <div className="profile-info">
        <p>Полное имя: {profile.fullName}</p>
        <p>Email: {profile.email}</p>
        <p>Дата рождения: {profile.birthDate}</p>
        <p>Место жительства: {profile.residence}</p>
        <p>Образование: {profile.education}</p>
        <p>Место рождения: {profile.birthPlace}</p>
        <p>Домашний адрес: {profile.homeAddress}</p>
      </div>

    </div>
  ) : (
    <p>Загрузка...</p>
  );
};

export default Profile;