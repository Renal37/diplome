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
      <h1 className="profile-header">Ваш профиль</h1>
      <div className="profile_menu">
        <div className="profile-avatar">
          <h1>Ваши данные</h1>
          <div className="profile-info">
            <p>Ваше имя: {profile.firstName}</p>
            <p>Email: {profile.email}</p>
            <p>Дата рождения: {profile.birthDate}</p>
            <p>Место жительства: {profile.residence}</p>
            <p>Образование: {profile.education}</p>
          </div>
        </div>
        <div className="profile_course">
          <h1>Ваши пройденные курсы</h1>
          <div className="course_end">

          </div>

        </div>
      </div>

    </div>


  ) : (
    <p>Загрузка...</p>
  );
};

export default Profile;