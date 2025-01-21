import { Routes, Route, Link } from 'react-router-dom';
import Header from './component/header/header.component';
import Main from './component/main/main.component';
import Main_top_center from './component/main_components/main_class/main/main_top_center/main_top_center.component.jsx';
import Teacher_top_center from './component/main_components/main_class/teacher/main_top_center/main_top_center.component.jsx';
import Organization_top_center from './component/main_components/main_class/organization/main_top_center/main_top_center.component.jsx';
import Professional_top_center from './component/main_components/main_class/professional/main_top_center/main_top_center.component.jsx';
import Promotion_top_center from './component/main_components/main_class/promotion/main_top_center/main_top_center.component.jsx';
import './style.css';
import AdminComponent from './component/admin/admin_component.jsx';
import Registration from './component/auth/registration/registration_component.jsx';
import Profile from './component/auth/profile/profile_component.jsx';
import Authorization from './component/auth/authorization/authorization_component.jsx';

const routes = [
  { path: "/", element: <Main_top_center /> },
  { path: "teacher", element: <Teacher_top_center /> },
  { path: "organization", element: <Organization_top_center /> },
  { path: "professional", element: <Professional_top_center /> },
  { path: "promotion", element: <Promotion_top_center /> }
];

const App = () => {
  return (
    <div>
      <Header />

      <Routes>
        <Route path="/" element={<Main />}>
          {routes.map(({ path, element }) => (
            <Route key={path} path={path} element={element} />
          ))}
        </Route>
        <Route path="/admin/*" element={<AdminComponent />} />
        <Route path="/registration" element={<Registration />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/authorization" element={<Authorization />} />
      </Routes>
    </div>
  );
}

export default App;