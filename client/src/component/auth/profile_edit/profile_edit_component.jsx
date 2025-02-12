import "./profile_edit_component.css";
import ProfileEditDocument from "./profile_edit_document/profile_edit_document";
import ProfileEditDate from "./profile_edit_date/profile_edit_date";
import { Link, Route, Routes } from "react-router-dom";

const Profile_Edit = () => {
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
    </div>
  );
};

export default Profile_Edit;