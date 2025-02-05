import { useState, useEffect } from "react";
import "./profile_component.css";
import CheckCourse from "../check_course/check_course_component";


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
            <p>Ваше ФИО: {profile.lastName && profile.firstName && profile.middleName ? `${profile.lastName} ${profile.firstName} ${profile.middleName}` : "Данные отсутствуют"}</p>
            <p>Email: {profile.email || "Данные отсутствуют"}</p>
            <p>Дата рождения: {profile.birthDate || "Данные отсутствуют"}</p>
            <p>Место жительства: {profile.residence || "Данные отсутствуют"}</p>
            <p>Образование: {profile.education || "Данные отсутствуют"}</p>
          </div>
        </div>
        <div className="profile_course">
          <div className="profile_course_top">
            <div className="profile_course_text">
              <h1>Курсы:</h1>
            </div>
          </div>


          <div className="course_end">
            <CheckCourse />
          </div>
        </div>
      </div>
    </div>

  ) : (
    <p>Загрузка...</p>
  );
};

export default Profile;