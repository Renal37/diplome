import { useEffect, useState } from "react";
import "./admin_check_profile.css";

const AdminCheckProfile = () => {
    const [users, setUsers] = useState([]);
    const [fioSearchQuery, setFioSearchQuery] = useState(""); // Поиск по ФИО
    const [loginSearchQuery, setLoginSearchQuery] = useState(""); // Поиск по логину
    const [showFullInfo, setShowFullInfo] = useState(false);

    useEffect(() => {
        fetch("http://localhost:5000/users")
            .then((res) => res.json())
            .then((data) => setUsers(data));
    }, []);

    // Фильтрация пользователей по ФИО и логину
    const filteredUsers = users.filter((user) => {
        const fullName = [
            user.lastname || "",
            user.firstname || "",
            user.middlename || ""
        ].join(" ").toLowerCase(); // Собираем ФИО в одну строку
        const login = user.username?.toLowerCase() || ""; // Логин пользователя

        // Проверяем, соответствует ли пользователь хотя бы одному из условий
        return (
            (!fioSearchQuery || fullName.includes(fioSearchQuery.toLowerCase())) && // Поиск по ФИО
            (!loginSearchQuery || login.includes(loginSearchQuery.toLowerCase())) // Поиск по логину
        );
    });

    return (
        <div className="admin-delete-component">
            <h1>Проверка профилей</h1>
            <div className="controls">
                <div className="search-container">
                    {/* Ввод для поиска по ФИО */}
                    <input
                        type="text"
                        placeholder="Поиск по ФИО"
                        value={fioSearchQuery}
                        onChange={(e) => setFioSearchQuery(e.target.value)}
                        className="search-input"
                    />
                    {/* Ввод для поиска по логину */}
                    <input
                        type="text"
                        placeholder="Поиск по логину"
                        value={loginSearchQuery}
                        onChange={(e) => setLoginSearchQuery(e.target.value)}
                        className="search-input"
                    />
                    <button
                        className="toggle-info-button"
                        onClick={() => setShowFullInfo(!showFullInfo)}
                    >
                        {showFullInfo ? "Скрыть детали" : "Показать детали"}
                    </button>
                </div>
            </div>
            <table className="course-table">
                <thead>
                    <tr>
                        <th>Логин</th>
                        <th>Email</th>
                        <th>ФИО</th>
                        {showFullInfo && (
                            <>
                                <th>Дата рождения</th>
                                <th>Место рождения</th>
                                <th>Образование</th>
                                <th>Место работы</th>
                                <th>Должность</th>
                                <th>Домашний адрес</th>
                                <th>Номер</th>
                                <th>Данные паспорта</th>
                                <th>СНИЛС</th>
                            </>
                        )}
                    </tr>
                </thead>
                <tbody>
                    {filteredUsers.length > 0 ? (
                        filteredUsers.map((user) => (
                            <tr key={user._id}>
                                <td>{user.username || "Данные отсутствуют"}</td>
                                <td>{user.email || "Данные отсутствуют"}</td>
                                <td>
                                    {user.lastname && user.firstname && user.middlename
                                        ? `${user.lastname} ${user.firstname} ${user.middlename}`
                                        : "Данные отсутствуют"}
                                </td>
                                {showFullInfo && (
                                    <>
                                        <td>{user.birthdate || "Данные отсутствуют"}</td>
                                        <td>{user.birthplace || "Данные отсутствуют"}</td>
                                        <td>{user.education || "Данные отсутствуют"}</td>
                                        <td>{user.workplace || "Данные отсутствуют"}</td>
                                        <td>{user.jobtitle || "Данные отсутствуют"}</td>
                                        <td>{user.homeaddress || "Данные отсутствуют"}</td>
                                        <td>{user.phone || "Данные отсутствуют"}</td>
                                        <td>{user.passportdata || "Данные отсутствуют"}</td>
                                        <td>{user.snils || "Данные отсутствуют"}</td>
                                    </>
                                )}
                            </tr>
                        ))
                    ) : (
                        <tr>
                            <td colSpan={showFullInfo ? 12 : 3} className="no-data-message">
                                Нет данных для отображения
                            </td>
                        </tr>
                    )}
                </tbody>
            </table>
        </div>
    );
};

export default AdminCheckProfile;