import "./profile_edit_component.css";
import ProfileEditDocument from "./profile_edit_document/profile_edit_document";
import ProfileEditDate from "./profile_edit_date/profile_edit_date";
import { Link, Route, Routes } from "react-router-dom";
import { useState, useEffect } from "react";

const Profile_Edit = () => {
  const [profile, setProfile] = useState(null);
 

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
      } 
    };

    fetchProfile();
  }, []);

  const handleDownloadContract = async () => {
    try {
      const response = await fetch(`http://localhost:5000/user/download-document`, {
        method: "GET",
        credentials: "include",
      });
      if (!response.ok) {
        throw new Error("Ошибка при скачивании договора");
      }
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.style.display = "none";
      a.href = url;
      a.download = "согласние.pdf";
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      alert("Вы скачали соглашение, теперь заполните его и отправьте нам!");
    } catch (err) {
      alert(err.message);
    }
  };

  const handleUploadContract = async (userId, file) => {
    const formData = new FormData();
    formData.append("contract", file);

    try {
      const response = await fetch(`http://localhost:5000/user/upload-document/${userId}`, {
        method: "POST",
        body: formData,
        credentials: "include",
      });
      if (!response.ok) {
        throw new Error("Ошибка при загрузке договора");
      }
      alert("Договор успешно загружен!");
    } catch (err) {
      alert(err.message);
    }
  };

  return (
    <div className="profile-container_edit">
      <h1 className="profile-header">Редактирование профиля</h1>
      <div className="profile_nav">
        <Link className="a" to="/auth/edit_profile">Ваши данные</Link>
        <Link className="a" to="document">Ваши документы</Link>
      </div>

      <div className="profile-info">
        <Routes>
          <Route index element={<ProfileEditDate />} />
          <Route path="document" element={<ProfileEditDocument />} />
        </Routes>
      </div>
      {profile && (
        <>
          <button
            className="download-contract-button"
            onClick={() => handleDownloadContract(profile.ID)}
          >
            Скачать договор
          </button>
          <label htmlFor="upload-contract" className="upload-label">
            Загрузить файл
          </label>
          <input
            id="upload-contract"
            type="file"
            accept=".pdf"
            style={{ display: "none" }}
            onChange={(e) => {
              const file = e.target.files[0];
              if (file) {
                handleUploadContract(profile.ID, file);
              }
            }}
          />
        </>
      )}
    </div>
  );
};

export default Profile_Edit;