import React, { useState, useEffect } from "react";
import "./admin_check_manage.css";

const AdminCoursesManagement = () => {
    const [registrations, setRegistrations] = useState([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState("");
    const [rejectReason, setRejectReason] = useState(""); // Причина отклонения
    const [expelReason, setExpelReason] = useState(""); // Причина отчисления
    const [selectedRejectRegistrationId, setSelectedRejectRegistrationId] = useState(null); // ID выбранной заявки для отклонения
    const [selectedExpelRegistrationId, setSelectedExpelRegistrationId] = useState(null); // ID выбранной заявки для отчисления
    const [selectedPdfRegistrationId, setSelectedPdfRegistrationId] = useState(null); // ID заявки с PDF
    const [pdfUrl, setPdfUrl] = useState(""); // URL для просмотра PDF
    const [selectedUser, setSelectedUser] = useState(null); // Выбранный пользователь
    const [filterStatus, setFilterStatus] = useState("all"); // Фильтр по статусу
    const [groups, setGroups] = useState([]); // Список групп
    const [searchUserName, setSearchUserName] = useState(""); // Поиск по имени пользователя
    const [searchCourseTitle, setSearchCourseTitle] = useState(""); // Поиск по названию курса

    // Фильтрация заявок по статусу, имени пользователя и названию курса
    const filteredRegistrations = registrations.filter((registration) => {
        const matchesStatus = filterStatus === "all" || registration.status === filterStatus;
        const matchesUserName = registration.userName.toLowerCase().includes(searchUserName.toLowerCase());
        const matchesCourseTitle = registration.courseTitle.toLowerCase().includes(searchCourseTitle.toLowerCase());
        return matchesStatus && matchesUserName && matchesCourseTitle;
    });

    // Загрузка списка групп
    useEffect(() => {
        const fetchGroups = async () => {
            try {
                const response = await fetch("http://localhost:5000/groups", {
                    credentials: "include",
                });
                if (!response.ok) {
                    throw new Error("Ошибка при загрузке групп");
                }
                const data = await response.json();
                setGroups(data.groups || []);
            } catch (error) {
                console.error("Error fetching groups:", error);
                setError("Ошибка при загрузке групп");
            }
        };
        fetchGroups();
    }, []);
    useEffect(() => {
        const handleKeyDown = (event) => {
            if (event.key === "Escape") {
                if (selectedRejectRegistrationId) {
                    setSelectedRejectRegistrationId(null);
                    setRejectReason("");
                }
                if (selectedExpelRegistrationId) {
                    setSelectedExpelRegistrationId(null);
                    setRejectReason("");
                }
                if (selectedExpelRegistrationId) {
                    setSelectedExpelRegistrationId(null);
                    setRejectReason("");
                }
                if (selectedPdfRegistrationId) {
                    setSelectedPdfRegistrationId(null);
                }
                if (selectedUser) {
                    setSelectedUser(null);
                }
            }
        };

        document.addEventListener("keydown", handleKeyDown);

        return () => {
            document.removeEventListener("keydown", handleKeyDown);
        };
    }, [selectedRejectRegistrationId, selectedExpelRegistrationId, selectedPdfRegistrationId, selectedUser]);

    useEffect(() => {
        console.log("Registrations updated:", registrations);
    }, [registrations]);

    // Загрузка данных при монтировании
    useEffect(() => {
        fetch("http://localhost:5000/admin/course-registrations", {
            credentials: "include",
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    setError(data.error);
                } else {
                    setRegistrations(data);
                }
                setIsLoading(false);
            })
            .catch((error) => {
                console.error("Error fetching registrations:", error);
                setError("Ошибка при загрузке заявок");
                setIsLoading(false);
            });
    }, []);

    // Одобрение заявки
    const handleApprove = (registrationId) => {
        fetch(`http://localhost:5000/admin/approve-registration/${registrationId}`, {
            method: "POST",
            credentials: "include",
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.success) {
                    setRegistrations(
                        registrations.map((reg) =>
                            reg._id === registrationId ? { ...reg, status: "Одобренный" } : reg
                        )
                    );
                } else {
                    setError(data.message || "Ошибка при одобрении заявки");
                }
            })
            .catch((error) => {
                console.error("Error approving registration:", error);
                setError("Ошибка при одобрении заявки");
            });
    };

    // Отклонение заявки
    const handleReject = (registrationId) => {
        setSelectedRejectRegistrationId(registrationId);
    };

    const confirmReject = () => {
        if (!rejectReason) {
            setError("Укажите причину отклонения");
            return;
        }
        fetch(`http://localhost:5000/admin/reject-registration/${selectedRejectRegistrationId}`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify({ reason: rejectReason }),
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.success) {
                    setRegistrations(
                        registrations.map((reg) =>
                            reg._id === selectedRejectRegistrationId
                                ? { ...reg, status: "Отклоненный", rejectReason }
                                : reg
                        )
                    );
                    setSelectedRejectRegistrationId(null);
                    setRejectReason("");
                } else {
                    setError(data.message || "Ошибка при отклонении заявки");
                }
            })
            .catch((error) => {
                console.error("Error rejecting registration:", error);
                setError("Ошибка при отклонении заявки");
            });
    };

    const handleViewConsent = (userId) => {
        window.open(`http://localhost:5000/user/view-consent/${userId}`, '_blank');
    };

    // Удаление заявки
    const handleDelete = (registrationId) => {
        fetch(`http://localhost:5000/admin/delete-registration/${registrationId}`, {
            method: "POST",
            credentials: "include",
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.success) {
                    setRegistrations(registrations.filter((reg) => reg._id !== registrationId));
                } else {
                    setError(data.message || "Ошибка при удалении заявки");
                }
            })
            .catch((error) => {
                console.error("Error deleting registration:", error);
                setError("Ошибка при удалении заявки");
            });
    };

    // Отчисление пользователя
    const handleExpel = (registrationId) => {
        setSelectedExpelRegistrationId(registrationId);
    };

    const confirmExpel = () => {
        if (!expelReason) {
            setError("Укажите причину отчисления");
            return;
        }
        fetch(`http://localhost:5000/admin/expel-registration/${selectedExpelRegistrationId}`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
            body: JSON.stringify({ reason: expelReason }),
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.success) {
                    setRegistrations(
                        registrations.map((reg) =>
                            reg._id === selectedExpelRegistrationId
                                ? { ...reg, status: "Отчисленный", rejectReason: expelReason }
                                : reg
                        )
                    );
                    setSelectedExpelRegistrationId(null);
                    setExpelReason("");
                } else {
                    setError(data.message || "Ошибка при отчислении");
                }
            })
            .catch((error) => {
                console.error("Error expelling registration:", error);
                setError("Ошибка при отчислении");
            });
    };

    // Просмотр договора (PDF)
    const handleViewPdf = (registrationId) => {
        setSelectedPdfRegistrationId(registrationId);
        setPdfUrl(`http://localhost:5000/user/view-contract/${registrationId}`);
    };
    const handleApprovePdf = (registrationId) => {
        fetch(`http://localhost:5000/admin/approve-contract/${registrationId}`, {
            method: "POST",
            credentials: "include",
        })
            .then((response) => response.json())
            .then((data) => {
                if (data.success) {
                    setRegistrations((prevRegistrations) =>
                        prevRegistrations.map((reg) =>
                            reg._id === registrationId ? { ...reg, status: "Принят" } : reg
                        )
                    );
                    setSelectedPdfRegistrationId(null);
                } else {
                    window.alert(data.message || "Ошибка при принятии заявки");
                }
            })
            .catch((error) => {
                console.error("Error approving contract:", error);
                window.alert("Произошла ошибка при принятии заявки. Пожалуйста, попробуйте снова.");
            });
    };


    const handleAssignGroup = async (registrationId, groupId) => {
        if (!groupId) {
            setError("Выберите группу");
            return;
        }

        try {
            const response = await fetch(`http://localhost:5000/admin/assign-group/${registrationId}`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ groupId }),
            });

            const data = await response.json();

            if (data.success) {
                setRegistrations((prevRegistrations) =>
                    prevRegistrations.map((reg) =>
                        reg._id === registrationId ? { ...reg, groupId } : reg
                    )
                );
            } else {
                setError(data.message || "Ошибка при привязке к группе");
            }
        } catch (error) {
            console.error("Error assigning group:", error);
            setError("Ошибка при привязке к группе");
        }
    };

    // Просмотр информации о пользователе
    const handleViewUser = (user) => {
        setSelectedUser(user);
    };

    const handleCloseUserInfo = () => {
        setSelectedUser(null);
    };

    if (isLoading) return <div>Loading...</div>;

    return (
        <div className="admin-approval-page">
            {error && <p className="error">{error}</p>}
            {/* Фильтр по статусу и поиск */}
            <div className="controls">
                <div className="search-container">
                    <input
                        type="text"
                        value={searchUserName}
                        onChange={(e) => setSearchUserName(e.target.value)}
                        placeholder="Поиск по имени пользователя"
                        className="search-input"
                    />
                    <input
                        type="text"
                        value={searchCourseTitle}
                        onChange={(e) => setSearchCourseTitle(e.target.value)}
                        placeholder="Поиск по названию курса"
                        className="search-input"
                    />
                    <div className="filter-section">
                        <label>Фильтр по статусу:</label>
                        <select value={filterStatus} onChange={(e) => setFilterStatus(e.target.value)}>
                            <option value="all">Все</option>
                            <option value="Ожидание">Ожидание</option>
                            <option value="Одобренный">Одобренный</option>
                            <option value="Отклоненный">Отклоненный</option>
                            <option value="Отчисленный">Отчисленный</option>
                            <option value="Принят">Принят</option>
                        </select>
                    </div>
                </div>
            </div>

            {/* Таблица заявок */}
            <table>
                <thead>
                    <tr>
                        <th>Курс</th>
                        <th>Пользователь</th>
                        <th>Группа</th>
                        <th>Статус</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    {filteredRegistrations.map((registration) => (
                        <tr key={registration._id}>
                            <td className="course-titles">{registration.courseTitle}</td>
                            <td>
                                <button onClick={() => handleViewUser({
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
                                })}>
                                    {registration.userName}
                                </button>
                            </td>
                            <td>
                                {registration.status === "Принят" && (
                                    <div className="filter-group">
                                        <select
                                            value={registration.groupId || ""}
                                            onChange={(e) => handleAssignGroup(registration._id, e.target.value)}
                                        >
                                            <option value="">Выберите группу</option>
                                            {groups.map((group) => (
                                                <option key={group._id} value={group._id}>
                                                    {group.groupName}
                                                </option>
                                            ))}
                                        </select>
                                    </div>
                                )}
                            </td>
                            <td>{registration.status}</td>
                            <td className="admin_btns">
                                {registration.status === "Ожидание" && (
                                    <>
                                        <button className="approve-btn" onClick={() => handleApprove(registration._id)}>Одобрить</button>
                                        <button className="reject-btn" onClick={() => handleReject(registration._id)}>Отклонить</button>
                                    </>
                                )}
                                {(registration.status === "Отклоненный" || registration.status === "Отчисленный") && (
                                    <button className="reject-btn" onClick={() => handleDelete(registration._id)}>Удалить</button>
                                )}
                                {registration.status === "Одобренный" && registration.contractFilePath && (
                                    <>
                                        <button className="reject-btn" onClick={() => handleReject(registration._id)}>Отклонить</button>
                                        <button className="approve-btn" onClick={() => handleViewPdf(registration._id)}>Проверить договор</button>
                                    </>
                                )}
                                {registration.status === "Одобренный" && !(registration.contractFilePath) && (
                                    <>
                                        <button className="reject-btn" onClick={() => handleReject(registration._id)}>Отклонить</button>
                                    </>
                                )}
                                {registration.status === "Принят" && registration.contractFilePath && (
                                    <button className="reject-btn" onClick={() => handleExpel(registration._id)}>Отчислить</button>
                                )}
                                {!(registration.status == "Ожидание") && !(registration.status === "Отклоненный" || registration.status === "Отчисленный") && (
                                    <>
                                        <button onClick={() => handleViewConsent(registration.userId)}>Просмотр согласия</button>
                                    </>
                                )}
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>

            {/* Модальное окно для отклонения заявки */}
            {selectedRejectRegistrationId && (
                <div className="modal">
                    <div className="modal-content">
                        <h3>Укажите причину отклонения</h3>
                        <textarea
                            type="text"
                            value={rejectReason}
                            onChange={(e) => setRejectReason(e.target.value)}
                            placeholder="Причина отклонения"
                        />
                        <button className="approve-btn" onClick={confirmReject}>Подтвердить</button>
                        <button className="reject-btn" onClick={() => setSelectedRejectRegistrationId(null)}>Отмена</button>
                    </div>
                </div>
            )}

            {/* Модальное окно для отчисления */}
            {selectedExpelRegistrationId && (
                <div className="modal">
                    <div className="modal-content">
                        <h3>Укажите причину отчисления</h3>
                        <textarea
                            type="text"
                            value={expelReason}
                            onChange={(e) => setExpelReason(e.target.value)}
                            placeholder="Причина отчисления"
                        />
                        <button className="approve-btn" onClick={confirmExpel}>Подтвердить</button>
                        <button className="reject-btn" onClick={() => setSelectedExpelRegistrationId(null)}>Отмена</button>
                    </div>
                </div>
            )}

            {/* Модальное окно для просмотра PDF */}
            {selectedPdfRegistrationId && (
                <div className="modal">
                    <div className="modal-content_for_view">
                        <h3>Проверка договора</h3>
                        <iframe
                            src={pdfUrl}
                            width="100%"
                            height="100%"
                            style={{ border: "none" }}
                        />
                        <div className="btns">
                            <button className="reject-btn" onClick={() => setSelectedPdfRegistrationId(null)}>Закрыть</button>
                            <button className="approve-btn" onClick={() => handleApprovePdf(selectedPdfRegistrationId)}>Принять</button>
                        </div>
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
                        <button className="reject-btn" onClick={handleCloseUserInfo}>Закрыть</button>
                    </div>
                </div>
            )}
        </div>
    );
};

export default AdminCoursesManagement;