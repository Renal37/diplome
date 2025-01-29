import { useState, useEffect } from "react";

const Profile = () => {
  const [user, setUser] = useState(null);
  const [password, setPassword] = useState("");
  const [education, setEducation] = useState("");
  const [passportFile, setPassportFile] = useState(null);
  const [snilsFile, setSnilsFile] = useState(null);

  useEffect(() => {
    fetch("http://localhost:5000/profile", {
      method: "GET",
      credentials: "include", // Включаем куки в запрос
    })
      .then(res => {
        if (!res.ok) {
          throw new Error("Ошибка при получении профиля");
        }
        return res.json();
      })
      .then(data => setUser(data))
      .catch(err => console.error(err));
  }, []);

  const updateProfile = () => {
    const formData = new FormData();
    formData.append("password", password);
    formData.append("education", education);
    if (passportFile) formData.append("passportPhoto", passportFile);
    if (snilsFile) formData.append("snilsPhoto", snilsFile);

    fetch(`http://localhost:5000/update-profile/${user._id}`, {
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

  return user ? (
    <div>
      <h1>Профиль {user.username}</h1>
      <p>Полное имя: {user.fullName}</p>
      <p>Email: {user.email}</p>
      <p>Дата рождения: {user.birthDate}</p>
      <p>Место жительства: {user.residence}</p>
      <p>Образование: {user.education}</p>
      <p>Место рождения: {user.birthPlace}</p>
      <p>Домашний адрес: {user.homeAddress}</p>
      <input type="password" placeholder="Новый пароль" onChange={e => setPassword(e.target.value)} />
      <select onChange={e => setEducation(e.target.value)}>
        <option value="Среднее">Среднее</option>
        <option value="Высшее">Высшее</option>
      </select>
      <input type="file" onChange={e => setPassportFile(e.target.files[0])} />
      <input type="file" onChange={e => setSnilsFile(e.target.files[0])} />
      <button onClick={updateProfile}>Обновить профиль</button>
    </div>
  ) : (
    <p>Загрузка...</p>
  );
};

export default Profile;