import { Routes, Route } from 'react-router-dom';
import Header from './component/header/header.component';
import Main from './component/main/main.component';
import Main_top_center from './component/main_components/main_class/main/main_top_center/main_top_center.component.jsx';
import Teacher_top_center from './component/main_components/main_class/teacher/main_top_center/main_top_center.component.jsx';
import Organization_top_center from './component/main_components/main_class/organization/main_top_center/main_top_center.component.jsx';
import Professional_top_center from './component/main_components/main_class/professional/main_top_center/main_top_center.component.jsx';
import Promotion_top_center from './component/main_components/main_class/promotion/main_top_center/main_top_center.component.jsx';
import './style.css';
import AdminComponent from './component/admin/admin_component.jsx';
import Auth from './component/auth/auth_component.jsx';
import Profile_Edit from './component/auth/profile_edit/profile_edit_component.jsx';
import CourseRegistration from './component/course/ course_component.jsx';
import Sveden from './component/sveden/sveden_component.jsx';
import Document from './component/document/document_component.jsx';

const routes = [
  { path: "/", element: <Main_top_center /> },
  { path: "teacher", element: <Teacher_top_center /> },
  { path: "organization", element: <Organization_top_center /> },
  { path: "professional", element: <Professional_top_center /> },
  { path: "promotion", element: <Promotion_top_center /> },
  { path: "sveden", element: <Sveden /> },
  { path: "document", element: <Document /> }
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
        <Route path="/auth/*" element={<Auth />} />
        <Route path="/auth/edit_profile/*" element={<Profile_Edit />} />
        <Route path="/courses/register/:courseId" element={<CourseRegistration/>} />
      </Routes>
    </div>
  );
}

export default App; 