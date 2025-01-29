import { useState, useEffect } from "react";
import "./profile_component.css";

const Profile = () => {
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [password, setPassword] = useState("");
  const [education, setEducation] = useState("");
  const [passportFile, setPassportFile] = useState(null);
  const [snilsFile, setSnilsFile] = useState(null);

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const response = await fetch("http://localhost:5000/profile", {
          method: "GET",
          credentials: "include", // ОБЯЗАТЕЛЬНО для отправки куки
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

  const updateProfile = () => {
    const formData = new FormData();
    formData.append("password", password);
    formData.append("education", education);
    if (passportFile) formData.append("passportPhoto", passportFile);
    if (snilsFile) formData.append("snilsPhoto", snilsFile);

    fetch(`http://localhost:5000/update-profile/${profile._id}`, {
      method: "PUT",
      credentials: "include", // Включаем куки в запрос
      body: formData,
    })
      .then(res => {
        if (!res.ok) {
          throw new Error("Ошибка при обновлении профиля");
        }
        return res.json();
      })
      .then(() => alert("Профиль обновлен"))
      .catch(err => console.error(err));
  };

  if (loading) {
    return <p>Загрузка...</p>;
  }

  if (error) {
    return <p>{error}</p>;
  }

  return profile ? (
    <div className="profile-container">
      <h1 className="profile-header">Профиль {profile.username}</h1>
      <div className="profile-info">
        <p>Полное имя: {profile.fullName}</p>
        <p>Email: {profile.email}</p>
        <p>Дата рождения: {profile.birthDate}</p>
        <p>Место жительства: {profile.residence}</p>
        <p>Образование: {profile.education}</p>
        <p>Место рождения: {profile.birthPlace}</p>
        <p>Домашний адрес: {profile.homeAddress}</p>
      </div>
      <input
        type="password"
        placeholder="Новый пароль"
        className="profile-input"
        onChange={e => setPassword(e.target.value)}
      />
      <select className="profile-input" onChange={e => setEducation(e.target.value)}>
        <option value="Среднее">Среднее</option>
        <option value="Высшее">Высшее</option>
      </select>
      <input
        type="file"
        className="profile-input"
        onChange={e => setPassportFile(e.target.files[0])}
      />
      <input
        type="file"
        className="profile-input"
        onChange={e => setSnilsFile(e.target.files[0])}
      />
      <button className="profile-button" onClick={updateProfile}>
        Обновить профиль
      </button>
    </div>
  ) : (
    <p>Загрузка...</p>
  );
};

export default Profile;