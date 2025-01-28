import { useState, useEffect } from "react";

const Profile = () => {
  const [user, setUser] = useState(null);
  const [password, setPassword] = useState("");
  const [education, setEducation] = useState("");
  const [passportFile, setPassportFile] = useState(null);
  const [snilsFile, setSnilsFile] = useState(null);

  useEffect(() => {
    fetch("http://localhost:5000/check-token", {
      method: "POST",
      headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
    }).then(res => res.json()).then(data => setUser(data));
  }, []);

  const updateProfile = () => {
    const formData = new FormData();
    formData.append("password", password);
    formData.append("education", education);
    if (passportFile) formData.append("passportPhoto", passportFile);
    if (snilsFile) formData.append("snilsPhoto", snilsFile);

    fetch(`http://localhost:5000/update-profile/${user._id}`, {
      method: "PUT",
      headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
      body: formData,
    }).then(() => alert("Профиль обновлен"));
  };

  return user ? (
    <div>
      <h1>Профиль {user.username}</h1>
      <input type="password" placeholder="Новый пароль" onChange={e => setPassword(e.target.value)} />
      <select onChange={e => setEducation(e.target.value)}>
        <option value="Среднее">Среднее</option>
        <option value="Высшее">Высшее</option>
      </select>
      <input type="file" onChange={e => setPassportFile(e.target.files[0])} />
      <input type="file" onChange={e => setSnilsFile(e.target.files[0])} />
      <button onClick={updateProfile}>Обновить профиль</button>
    </div>
  ) : (<p>Загрузка...</p>);
};

export default Profile; 