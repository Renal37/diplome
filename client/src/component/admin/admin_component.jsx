import { Route, Routes, Link } from 'react-router-dom';
import AdminDeleteComponent from '../admin_components/admin_delete/admin_delete_component';
import AdminUpdateComponent from '../admin_components/admin_update/admin_update_component';
import AdminAddComponent from '../admin_components/admin_add/admin_add_component';
import './admin_component.css';
import AdminApprovalPage from '../admin_components/admin_approval_page/admin_approval_page';
import AdminCheckProfile from '../admin_components/admin_check_profile/admin_check_profile';

const AdminComponent = () => {
  return (
    <div className="admin-component">
      <nav>
        <ul className="nav-links">
          <li>
            <Link to="/admin" className="nav-link">Добавить курс</Link>
          </li>
          <li>
            <Link to="/admin/update" className="nav-link">Обновить курс</Link>
          </li>
          <li>
            <Link to="/admin/delete" className="nav-link">Удалить курс</Link>
          </li>
          <li>
            <Link to="/admin/profile" className="nav-link">Просмотр курса</Link>
          </li>
          <li>
            <Link to="/admin/approval" className="nav-link">Одобрение курса</Link>
          </li>
        </ul>
      </nav>
      <div className="admin_components">
        <Routes>
          <Route index element={<AdminAddComponent />} />
          <Route path="update" element={<AdminUpdateComponent />} />
          <Route path="delete" element={<AdminDeleteComponent />} />
          <Route path="profile" element={<AdminCheckProfile />} />
          <Route path="approval" element={<AdminApprovalPage />} />
        </Routes>
      </div>

    </div>
  );
};

export default AdminComponent;