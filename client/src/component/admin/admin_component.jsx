import { Route, Routes, Link } from 'react-router-dom';
import './admin_component.css';
import AdminCheckProfile from '../admin_components/admin_check_profile/admin_check_profile';
import AdminAcceptPage from '../admin_components/admin_accept_page/admin_accept_page';
import AdminEditGroup from '../admin_components/admin_edit_group/admin_edit_group';
import AdminCreateGroup from '../admin_components/admin_create_group/admin_create_group';
import AdminCourseManagement from '../admin_components/admin_course_management/admin_course_management';
import AdminCoursesManagement from '../admin_components/admin_course_courses/admin_check_manage';

const AdminComponent = () => {
  return (
    <div className="admin-component container">
      <div>
        <h1 className="admin-title">Административная панель</h1>

      </div>
      <nav>
        <ul className="nav-links">
          <li>
            <Link to="/admin" className="nav-link">Управление курс</Link>
          </li>
       
          <li>
            <Link to="/admin/profile" className="nav-link">Просмотр всех профилей</Link>
          </li>
          <li>
            <Link to="/admin/courses" className="nav-link">Управление заявками</Link>
          </li>

          {/* <li>
            <Link to="/admin/accept" className="nav-link">Выдача докумета</Link>
          </li> */}
          <li>
            <Link to="/admin/creategroups" className="nav-link">Создание группы</Link>
          </li>
          <li>
            <Link to="/admin/editgroups" className="nav-link">Редакторивание группы</Link>
          </li>
        </ul>
      </nav>
      <div className="admin_components">
        <Routes>
          <Route index element={<AdminCourseManagement />} />
          <Route path="profile" element={<AdminCheckProfile />} />
          <Route path="courses" element={<AdminCoursesManagement />} />
          <Route path="accept" element={<AdminAcceptPage />} />
          <Route path="creategroups" element={<AdminCreateGroup />} />
          <Route path="editgroups" element={<AdminEditGroup />} />
        </Routes>
      </div>

    </div>
  );
};

export default AdminComponent;