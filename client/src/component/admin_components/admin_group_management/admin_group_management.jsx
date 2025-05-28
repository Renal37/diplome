import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import "./admin_group_management.css";

const AdminGroupManagement = () => {
    const navigate = useNavigate();
    const [groups, setGroups] = useState([]);
    const [filteredGroups, setFilteredGroups] = useState([]);
    const [searchTerm, setSearchTerm] = useState("");
    const [selectedGroup, setSelectedGroup] = useState(null);
    const [groupName, setGroupName] = useState("");
    const [courseId, setCourseId] = useState("");
    const [error, setError] = useState("");
    const [successMessage, setSuccessMessage] = useState("");
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [groupMembers, setGroupMembers] = useState([]);
    const [selectedGroupId, setSelectedGroupId] = useState(null);
    const [courses, setCourses] = useState([]);
    const [groupsWithMembers, setGroupsWithMembers] = useState({});
    const [orders, setOrders] = useState([]);
    const [selectedOrderId, setSelectedOrderId] = useState(null);
    const [isEnrollModalOpen, setIsEnrollModalOpen] = useState(false);
    const [isExpelModalOpen, setIsExpelModalOpen] = useState(false);
    const [rejectReason, setRejectReason] = useState("");
    const [isRejectModalOpen, setIsRejectModalOpen] = useState(false);

    // Загрузка групп
    const fetchGroups = async () => {
        try {
            const response = await fetch("http://localhost:5000/groups");
            if (!response.ok) throw new Error("Ошибка при загрузке групп");
            const data = await response.json();
            if (!data?.groups) throw new Error("Данные не получены или пусты");

            const groupsWithMembersInfo = {};
            for (const group of data.groups) {
                const membersResponse = await fetch(`http://localhost:5000/admin/group-members/${group._id}`);
                const membersData = await membersResponse.json();
                groupsWithMembersInfo[group._id] = membersData.members?.length > 0;
            }

            setGroups(data.groups);
            setFilteredGroups(data.groups);
            setGroupsWithMembers(groupsWithMembersInfo);
        } catch (error) {
            console.error("Error fetching groups:", error);
            setError("Ошибка при загрузке групп");
        }
    };

    // Загрузка курсов
    const fetchCourses = async () => {
        try {
            const response = await fetch("http://localhost:5000/courses");
            if (!response.ok) throw new Error("Ошибка при загрузке курсов");
            const data = await response.json();
            setCourses(data);
        } catch (error) {
            console.error("Error fetching courses:", error);
            setError("Ошибка при загрузке курсов");
        }
    };

    // Загрузка приказов
    const fetchOrders = async () => {
        try {
            const response = await fetch("http://localhost:5000/admin/orders", { credentials: "include" });
            if (!response.ok) throw new Error("Ошибка при загрузке приказов");
            const data = await response.json();
            setOrders(data || []);
        } catch (error) {
            console.error("Error fetching orders:", error);
            setError("Ошибка при загрузке приказов");
        }
    };

    useEffect(() => {
        fetchGroups();
        fetchCourses();
        fetchOrders();
    }, []);

    // Поиск
    const handleSearchChange = (e) => {
        const term = e.target.value.toLowerCase();
        setSearchTerm(term);
        setFilteredGroups(groups.filter(group =>
            group.groupName.toLowerCase().includes(term)
        ));
    };

    // Редактирование группы
    const handleEdit = (group) => {
        setSelectedGroup(group);
        setGroupName(group.groupName);
        setCourseId(group.courseId);
    };

    // Сохранение изменений
    const handleSave = async () => {
        if (!groupName || !courseId) {
            setError("Пожалуйста, заполните название группы и выберите курс");
            return;
        }
        try {
            const response = await fetch(`http://localhost:5000/admin/update-group/${selectedGroup._id}`, {
                method: "PUT",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ groupName, courseId })
            });
            const data = await response.json();
            if (data.success) {
                setSuccessMessage("Группа успешно обновлена!");
                handleClose();
                fetchGroups();
            } else {
                setError(data.message || "Ошибка при обновлении группы");
            }
        } catch (error) {
            console.error("Error updating group:", error);
            setError("Ошибка при обновлении группы");
        }
    };

    // Удаление группы
    const handleDelete = async (groupId) => {
        try {
            const response = await fetch(`http://localhost:5000/admin/delete-group/${groupId}`, {
                method: "DELETE",
                headers: { "Content-Type": "application/json" },
                credentials: "include"
            });
            const data = await response.json();
            if (data.success) {
                setSuccessMessage("Группа успешно удалена!");
                fetchGroups();
            } else {
                setError(data.message || "Ошибка при удалении группы");
            }
        } catch (error) {
            console.error("Error deleting group:", error);
            setError("Ошибка при удалении группы");
        }
    };

    // Отклонение записей группы
    const handleRejectGroup = (groupId) => {
        setSelectedGroupId(groupId);
        setIsRejectModalOpen(true);
    };

    const confirmReject = async () => {
        if (!rejectReason) {
            setError("Укажите причину отклонения");
            return;
        }
        try {
            const response = await fetch(`http://localhost:5000/admin/reject-group-registrations/${selectedGroupId}`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ reason: rejectReason })
            });
            const data = await response.json();
            if (data.success) {
                setSuccessMessage("Записи группы успешно отклонены!");
                setIsRejectModalOpen(false);
                setRejectReason("");
                fetchGroups();
            } else {
                setError(data.message || "Ошибка при отклонении записей");
            }
        } catch (error) {
            console.error("Error rejecting group registrations:", error);
            setError("Ошибка при отклонении записей");
        }
    };

    // Создание группы
    const handleCreateGroup = async (e) => {
        e.preventDefault();
        if (!groupName || !courseId) {
            setError("Пожалуйста, заполните название группы и выберите курс");
            return;
        }
        try {
            const response = await fetch("http://localhost:5000/admin/create-group", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ groupName, courseId })
            });
            const data = await response.json();
            if (data.success) {
                setSuccessMessage("Группа успешно создана!");
                handleClose();
                fetchGroups();
            } else {
                setError(data.message || "Ошибка при создании группы");
            }
        } catch (error) {
            console.error("Error creating group:", error);
            setError("Ошибка при создании группы");
        }
    };

    // Получение участников
    const fetchGroupMembers = async (groupId) => {
        try {
            const response = await fetch(`http://localhost:5000/admin/group-members/${groupId}`);
            if (!response.ok) throw new Error("Ошибка при загрузке участников");
            const data = await response.json();
            setGroupMembers(data.members || []);
        } catch (error) {
            console.error("Error fetching members:", error);
            setError("Ошибка при загрузке участников");
        }
    };

    // Просмотр участников
    const handleViewMembers = (groupId) => {
        setSelectedGroupId(groupId);
        fetchGroupMembers(groupId);
    };

    // Зачисление группы
    const handleEnrollGroup = (groupId) => {
        setSelectedGroupId(groupId);
        setIsEnrollModalOpen(true);
    };

    const confirmEnroll = async () => {
        if (!selectedOrderId) {
            setError("Выберите приказ о зачислении");
            return;
        }
        try {
            const response = await fetch(`http://localhost:5000/admin/enroll-group/${selectedGroupId}`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ orderId: selectedOrderId })
            });
            const data = await response.json();
            if (data.success) {
                setSuccessMessage("Группа успешно зачислена!");
                setIsEnrollModalOpen(false);
                setSelectedOrderId(null);
                fetchGroupMembers(selectedGroupId);
            } else {
                setError(data.message || "Ошибка при зачислении группы");
            }
        } catch (error) {
            console.error("Error enrolling group:", error);
            setError("Ошибка при зачислении группы");
        }
    };

    // Отчисление/выпуск группы
    const handleExpelGroup = (groupId) => {
        setSelectedGroupId(groupId);
        setIsExpelModalOpen(true);
    };

    const confirmExpel = async () => {
        if (!selectedOrderId) {
            setError("Выберите приказ");
            return;
        }
        try {
            const response = await fetch(`http://localhost:5000/admin/expel-group/${selectedGroupId}`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ orderId: selectedOrderId })
            });
            const data = await response.json();
            if (data.success) {
                setSuccessMessage("Приказ успешно выдан!");
                setIsExpelModalOpen(false);
                setSelectedOrderId(null);
                fetchGroupMembers(selectedGroupId);
            } else {
                setError(data.message || "Ошибка при выдаче приказа");
            }
        } catch (error) {
            console.error("Error expelling group:", error);
            setError("Ошибка при выдаче приказа");
        }
    };

    // Закрытие модальных окон
    const handleCloseMembersModal = () => {
        setSelectedGroupId(null);
        setGroupMembers([]);
    };

    const handleClose = () => {
        setIsModalOpen(false);
        setSelectedGroup(null);
        setGroupName("");
        setCourseId("");
        setError("");
        setSuccessMessage("");
    };

    useEffect(() => {
        const handleKeyDown = (e) => {
            if (e.key === "Escape") {
                handleClose();
                handleCloseMembersModal();
                setIsEnrollModalOpen(false);
                setIsExpelModalOpen(false);
                setIsRejectModalOpen(false);
                setSelectedOrderId(null);
                setRejectReason("");
            }
        };
        window.addEventListener("keydown", handleKeyDown);
        return () => window.removeEventListener("keydown", handleKeyDown);
    }, []);

    return (
        <div className="admin-groups-management-page">
            {error && <div className="error-message">{error}</div>}
            {successMessage && <div className="success-message">{successMessage}</div>}

            <div className="admin-course-header">
                <input
                    type="text"
                    placeholder="Поиск по группам..."
                    value={searchTerm}
                    onChange={handleSearchChange}
                    className="search-input"
                />
                <button className="approve-btn" onClick={() => setIsModalOpen(true)}>
                    Создать группу
                </button>
            </div>

            {isModalOpen && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Создание новой группы</h2>
                        <form onSubmit={handleCreateGroup}>
                            <label>Название группы:</label>
                            <input
                                type="text"
                                value={groupName}
                                onChange={(e) => setGroupName(e.target.value)}
                                placeholder="Введите название группы"
                            />
                            <label>Курс:</label>
                            <select
                                value={courseId}
                                onChange={(e) => setCourseId(e.target.value)}
                            >
                                <option value="">Выберите курс</option>
                                {courses.map(course => (
                                    <option key={course._id} value={course._id}>
                                        {course.title} (Макс. студентов: {course.maxStudents})
                                    </option>
                                ))}
                            </select>
                            <div className="form-buttons">
                                <button className="approve-btn" type="submit">Создать</button>
                                <button className="reject-btn" type="button" onClick={handleClose}>Закрыть</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            <table>
                <thead>
                    <tr>
                        <th>Название группы</th>
                        <th>Курс</th>
                        <th>Макс. студентов</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    {filteredGroups.map((group) => (
                        <tr key={group._id}>
                            <td>{group.groupName}</td>
                            <td>
                                {courses.find((course) => course._id === group.courseId)?.title || "Неизвестный курс"}
                            </td>
                            <td>
                                {courses.find((course) => course._id === group.courseId)?.maxStudents || "Не указано"}
                            </td>
                            <td>
                                <div className="bt">
                                    {groupsWithMembers[group._id] && (
                                        <>
                                            <button className="approve-btn" onClick={() => handleViewMembers(group._id)}>
                                                Просмотр участников
                                            </button>
                                            <button className="approve-btn" onClick={() => handleEnrollGroup(group._id)}>
                                                Зачислить группу
                                            </button>
                                            <button className="approve-btn" onClick={() => handleExpelGroup(group._id)}>
                                                Выдать приказ
                                            </button>
                                        </>
                                    )}
                                    {!groupsWithMembers[group._id] && (
                                        <>
                                            <button className="approve-btn" onClick={() => handleEdit(group)}>
                                                Редактировать
                                            </button>
                                            <button className="reject-btn" onClick={() => handleRejectGroup(group._id)}>
                                                Отклонить записи
                                            </button>
                                            <button className="reject-btn" onClick={() => handleDelete(group._id)}>
                                                Удалить
                                            </button>
                                        </>
                                    )}
                                </div>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>

            {selectedGroup && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Редактирование группы</h2>
                        <form>
                            <label>Название группы:</label>
                            <input
                                type="text"
                                value={groupName}
                                onChange={(e) => setGroupName(e.target.value)}
                            />
                            <label>Курс:</label>
                            <select
                                value={courseId}
                                onChange={(e) => setCourseId(e.target.value)}
                            >
                                <option value="">Выберите курс</option>
                                {courses.map(course => (
                                    <option key={course._id} value={course._id}>
                                        {course.title} (Макс. студентов: {course.maxStudents})
                                    </option>
                                ))}
                            </select>
                            <div className="form-buttons">
                                <button type="button" onClick={handleSave}>Сохранить</button>
                                <button type="button" onClick={handleClose}>Закрыть</button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            {selectedGroupId && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Участники группы</h2>
                        <table>
                            <thead>
                                <tr>
                                    <th>Имя пользователя</th>
                                    <th>Email</th>
                                    <th>Статус</th>
                                    <th>Действия</th>
                                </tr>
                            </thead>
                            <tbody>
                                {groupMembers.map(member => (
                                    <tr key={member._id}>
                                        <td>{member.username}</td>
                                        <td>{member.email}</td>
                                        <td>{member.status || "Неизвестно"}</td>
                                        <td>
                                            <button
                                                className="toggle-info-button"
                                                onClick={() => navigate(`/admin/profile?username=${member.username}`)}
                                            >
                                                Посмотреть данные
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                        <div className="button">
                            <button className="reject-btn" onClick={handleCloseMembersModal}>
                                Закрыть
                            </button>
                        </div>
                    </div>
                </div>
            )}

            {isEnrollModalOpen && (
                <div className="modal">
                    <div className="modal-content">
                        <h3>Зачисление группы</h3>
                        <div className="form-group">
                            <label>Приказ о зачислении:</label>
                            <select
                                value={selectedOrderId || ""}
                                onChange={(e) => setSelectedOrderId(e.target.value)}
                            >
                                <option value="">Выберите приказ</option>
                                {orders
                                    .filter((order) => order.orderType === "О зачислении обучающихся")
                                    .map((order) => (
                                        <option key={order.id} value={order.id}>
                                            {order.number} от {new Date(order.date).toLocaleDateString()} ({order.orderType})
                                        </option>
                                    ))}
                            </select>
                        </div>
                        <button className="approve-btn" onClick={confirmEnroll}>
                            Подтвердить
                        </button>
                        <button
                            className="reject-btn"
                            onClick={() => {
                                setIsEnrollModalOpen(false);
                                setSelectedOrderId(null);
                            }}
                        >
                            Отмена
                        </button>
                    </div>
                </div>
            )}

            {isExpelModalOpen && (
                <div className="modal">
                    <div className="modal-content">
                        <h3>Выбор приказа</h3>
                        <div className="form-group">
                            <label>Приказ:</label>
                            <select
                                value={selectedOrderId || ""}
                                onChange={(e) => setSelectedOrderId(e.target.value)}
                            >
                                <option value="">Выберите приказ</option>
                                {orders
                                    .filter(
                                        (order) =>
                                            order.orderType === "О выпуске обучающихся" ||
                                            order.orderType === "Об отчислении обучающихся"
                                    )
                                    .map((order) => (
                                        <option key={order.id} value={order.id}>
                                            {order.number} от {new Date(order.date).toLocaleDateString()} ({order.orderType})
                                        </option>
                                    ))}
                            </select>
                        </div>
                        <button className="approve-btn" onClick={confirmExpel}>
                            Подтвердить
                        </button>
                        <button
                            className="reject-btn"
                            onClick={() => {
                                setIsExpelModalOpen(false);
                                setSelectedOrderId(null);
                            }}
                        >
                            Отмена
                        </button>
                    </div>
                </div>
            )}

            {isRejectModalOpen && (
                <div className="modal">
                    <div className="modal-content">
                        <h3>Укажите причину отклонения</h3>
                        <textarea
                            value={rejectReason}
                            onChange={(e) => setRejectReason(e.target.value)}
                            placeholder="Причина отклонения (например, мало людей)"
                        />
                        <button className="approve-btn" onClick={confirmReject}>
                            Подтвердить
                        </button>
                        <button
                            className="reject-btn"
                            onClick={() => {
                                setIsRejectModalOpen(false);
                                setRejectReason("");
                            }}
                        >
                            Отмена
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
};

export default AdminGroupManagement;