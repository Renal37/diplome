import React, { useState, useEffect } from "react";
import "./admin_accept_page.css";

const AdminAcceptPage = () => {
    const [registrations, setRegistrations] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState("");
    const [rejectReason, setRejectReason] = useState(""); // Причина отчисления
    const [selectedRegistrationId, setSelectedRegistrationId] = useState(null);
    const [selectedDocumentId, setSelectedDocumentId] = useState(null); // ID выбранной заявки
    const [selectedUser, setSelectedUser] = useState(null); // Выбранный пользователь
    const [documentType, setDocumentType] = useState(""); // Тип документа (сертификат/диплом)

    useEffect(() => {
        fetch("http://localhost:5000/admin/course-registrations", {
            credentials: "include",
        })
            .then((response) => response.json())
            .then((data) => {
                console.log("Data from backend:", data);
                if (data.error) {
                    setError(data.error);
                } else {
                    // Фильтруем заявки по статусу "Одобренный"
                    const approvedRegistrations = data.filter(
                        (reg) => reg.status === "Одобренный"
                    );
                    setRegistrations(approvedRegistrations);
                }
                setIsLoading(false);
            })
            .catch((error) => {
                console.error("Error fetching registrations:", error);
                setError("Ошибка при загрузке заявок");
                setIsLoading(false);
            });
    }, []);

    const handleExpel = (registrationId) => {
        setSelectedRegistrationId(registrationId);
    };

    const confirmExpel = () => {
        if (!rejectReason) {
            setError("Укажите причину отчисления");
            return;
        }

        fetch(`http://localhost:5000/admin/expel-registration/${selectedRegistrationId}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            credentials: "include",
            body: JSON.stringify({ reason: rejectReason }), // Отправляем причину отчисления
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.success) {
                    setRegistrations(
                        registrations.map((reg) =>
                            reg._id === selectedRegistrationId
                                ? { ...reg, status: "Отчисленный", rejectReason: rejectReason }
                                : reg
                        )
                    );
                    setRejectReason(""); // Очищаем поле причины
                    setSelectedRegistrationId(null); // Закрываем модальное окно
                } else {
                    setError(data.message || "Ошибка при отчислении");
                }
            })
            .catch((error) => {
                console.error("Error expelling registration:", error);
                setError("Ошибка при отчислении");
            });
    };

    const handleIssueDocument = (registrationId) => {
        setSelectedDocumentId(registrationId);
    };

    const confirmIssueDocument = () => {
        if (!documentType) {
            setError("Выберите тип документа");
            return;
        }

        fetch(`http://localhost:5000/admin/issue-document/${selectedRegistrationId}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            credentials: "include",
            body: JSON.stringify({ documentType }), // Отправляем тип документа
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.success) {
                    setRegistrations(
                        registrations.map((reg) =>
                            reg._id === selectedRegistrationId
                                ? { ...reg, documentType: documentType }
                                : reg
                        )
                    );
                    setDocumentType(""); // Очищаем выбор типа документа
                    setSelectedDocumentId(null); // Закрываем модальное окно
                } else {
                    setError(data.message || "Ошибка при выдаче документа");
                }
            })
            .catch((error) => {
                console.error("Error issuing document:", error);
                setError("Ошибка при выдаче документа");
            });
    };

    const handleViewUser = (user) => {
        setSelectedUser(user);
    };

    const handleCloseUserInfo = () => {
        setSelectedUser(null);
    };

    if (isLoading) {
        return <div>Загрузка...</div>;
    }

    if (error) {
        return <div>{error}</div>;
    }

    return (
        <div className="admin-approval-page">
            <h1>Заявки на курсы (Одобренные)</h1>
            <table>
                <thead>
                    <tr>
                        <th>Курс</th>
                        <th>Пользователь</th>
                        <th>Статус</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    {registrations.map((registration) => (
                        <tr key={registration._id}>
                            <td>{registration.courseTitle}</td>
                            <td>
                                <span
                                    style={{
                                        color: "#007bff",
                                        cursor: "pointer",
                                        textDecoration: "underline",
                                    }}
                                    onClick={() =>
                                        handleViewUser({
                                            username: registration.userName,
                                            email: registration.userEmail,
                                            fullname: registration.userFullname,
                                            birthdate: registration.userBirthdate,
                                            birthplace: registration.userBirthplace,
                                            education: registration.userEducation,
                                            workplace: registration.userWorkplace,
                                            jobtitle: registration.userJobtitle,
                                            homeaddress: registration.userHomeaddress,
                                            phone: registration.userPhone,
                                            passportdata: registration.userPassportdata,
                                            snils: registration.userSnils,
                                        })
                                    }
                                >
                                    {registration.userName}
                                </span>
                            </td>
                            <td>{registration.status}</td>
                            <td>
                                <button className="expel-btn" onClick={() => handleExpel(registration._id)}>
                                    Отчислить
                                </button>
                                <button
                                    className="issue-document-btn"
                                    onClick={() => handleIssueDocument(registration._id)}
                                >
                                    Выдать документ
                                </button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>

            {/* Модальное окно для отчисления */}
            {selectedRegistrationId && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Укажите причину отчисления</h2>
                        <textarea
                            value={rejectReason}
                            onChange={(e) => setRejectReason(e.target.value)}
                            placeholder="Причина отчисления"
                        />
                        <button onClick={confirmExpel}>Подтвердить</button>
                        <button onClick={() => setSelectedRegistrationId(null)}>Отмена</button>
                    </div>
                </div>
            )}

            {/* Модальное окно для выбора типа документа */}
            {selectedDocumentId && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Выберите тип документа</h2>
                        <select
                            value={documentType}
                            onChange={(e) => setDocumentType(e.target.value)}
                        >
                            <option value="">Выберите тип</option>
                            <option value="Сертификат">Сертификат</option>
                            <option value="Диплом">Диплом</option>
                        </select>
                        <button onClick={confirmIssueDocument}>Подтвердить</button>
                        <button onClick={() => setSelectedDocumentId(null)}>Отмена</button>
                    </div>
                </div>
            )}

            {/* Модальное окно для просмотра информации о пользователе */}
            {selectedUser && (
                <div className="user-info-modal">
                    <div className="user-info-container">
                        <h3>Информация о пользователе</h3>
                        <table className="user-info-table">
                            <tbody>
                                <tr>
                                    <th>Логин</th>
                                    <td>{selectedUser.username || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Email</th>
                                    <td>{selectedUser.email || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>ФИО</th>
                                    <td>{selectedUser.fullname || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Дата рождения</th>
                                    <td>{selectedUser.birthdate || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Место рождения</th>
                                    <td>{selectedUser.birthplace || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Образование</th>
                                    <td>{selectedUser.education || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Место работы</th>
                                    <td>{selectedUser.workplace || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Должность</th>
                                    <td>{selectedUser.jobtitle || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Домашний адрес</th>
                                    <td>{selectedUser.homeaddress || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Номер телефона</th>
                                    <td>{selectedUser.phone || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>Данные паспорта</th>
                                    <td>{selectedUser.passportdata || "Данные отсутствуют"}</td>
                                </tr>
                                <tr>
                                    <th>СНИЛС</th>
                                    <td>{selectedUser.snils || "Данные отсутствуют"}</td>
                                </tr>
                            </tbody>
                        </table>
                        <button onClick={handleCloseUserInfo}>Закрыть</button>
                    </div>
                </div>
            )}
        </div>
    );
};

export default AdminAcceptPage;