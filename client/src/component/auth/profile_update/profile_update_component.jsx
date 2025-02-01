import { useState, useEffect } from "react";
import "./profile_update_component.css";

const Profile_Edit = () => {
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [oldPassword, setOldPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [birthDate, setBirthDate] = useState("");
  const [residence, setResidence] = useState("");
  const [education, setEducation] = useState("");
  const [birthPlace, setBirthPlace] = useState("");
  const [homeAddress, setHomeAddress] = useState("");
  const [consent, setConsent] = useState(false);

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
        setFullName(data.fullName);
        setEmail(data.email);
        setBirthDate(data.birthDate);
        setResidence(data.residence);
        setEducation(data.education);
        setBirthPlace(data.birthPlace);
        setHomeAddress(data.homeAddress);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchProfile();
  }, []);

  const updateProfile = async () => {
    if (!consent) {
      alert("Необходимо согласие на обработку данных");
      return;
    }

    const formData = new FormData();
    formData.append("oldPassword", oldPassword);
    formData.append("newPassword", newPassword);
    formData.append("fullName", fullName);
    formData.append("email", email);
    formData.append("birthDate", birthDate);
    formData.append("residence", residence);
    formData.append("education", education);
    formData.append("birthPlace", birthPlace);
    formData.append("homeAddress", homeAddress);

    try {
      const response = await fetch(`http://localhost:5000/update-profile/${profile._id}`, {
        method: "PUT",
        credentials: "include",
        body: formData,
      });

      if (!response.ok) {
        throw new Error("Ошибка при обновлении профиля");
      }

      alert("Профиль успешно обновлен");
    } catch (err) {
      console.error(err);
      alert("Ошибка при обновлении профиля");
    }
  };

  if (loading) {
    return <p>Загрузка...</p>;
  }

  if (error) {
    return <p>{error}</p>;
  }

  return profile ? (
    <div className="profile-container">
      <h1 className="profile-header">Редактирование профиля {profile.username}</h1>
      <div className="profile-info">
        <label>
          ФИО:
          <input
            type="text"
            value={fullName}
            onChange={(e) => setFullName(e.target.value)}
          />
        </label>
        <label>
          Email:
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
        </label>
        <label>
          Дата рождения:
          <input
            type="date"
            value={birthDate}
            onChange={(e) => setBirthDate(e.target.value)}
          />
        </label>
        <label>
          Место жительства:
          <input
            type="text"
            value={residence}
            onChange={(e) => setResidence(e.target.value)}
          />
        </label>
        <label>
          Образование:
          <select value={education} onChange={(e) => setEducation(e.target.value)}>
            <option value="Среднее">Среднее</option>
            <option value="Высшее">Высшее</option>
          </select>
        </label>
        <label>
          Место рождения:
          <input
            type="text"
            value={birthPlace}
            onChange={(e) => setBirthPlace(e.target.value)}
          />
        </label>
        <label>
          Домашний адрес:
          <input
            type="text"
            value={homeAddress}
            onChange={(e) => setHomeAddress(e.target.value)}
          />
        </label>
        <label>
          Старый пароль:
          <input
            type={showPassword ? "text" : "password"}
            value={oldPassword}
            onChange={(e) => setOldPassword(e.target.value)}
          />
        </label>
        <label>
          Новый пароль:
          <input
            type={showPassword ? "text" : "password"}
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
          />
        </label>
        <label>
          <input
            type="checkbox"
            checked={showPassword}
            onChange={() => setShowPassword(!showPassword)}
          />
          Показать пароль
        </label>
        <label>
          <input
            type="checkbox"
            checked={consent}
            onChange={() => setConsent(!consent)}
          />
          Согласие на обработку данных
        </label>
      </div>
      <button className="profile-button" onClick={updateProfile}>
        Обновить профиль
      </button>
    </div>
  ) : (
    <p>Загрузка...</p>
  );
};

export default Profile_Edit;