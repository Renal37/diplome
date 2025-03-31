import React, { useState, useEffect } from 'react';
import './admin_order.css';

const AdminOrders = () => {
    const [orderTypes, setOrderTypes] = useState([]);
    const [orders, setOrders] = useState([]);
    const [error, setError] = useState('');
    const [newOrder, setNewOrder] = useState({
        number: '',
        date: new Date().toISOString().slice(0, 10), // Формат YYYY-MM-DD
        orderTypeId: '',
        description: ''
    });
    const [newOrderType, setNewOrderType] = useState('');
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [typesRes, ordersRes] = await Promise.all([
                    fetch('http://localhost:5000/admin/order-types', { credentials: 'include' }),
                    fetch('http://localhost:5000/admin/orders', { credentials: 'include' })
                ]);

                const typesData = await typesRes.json();
                const ordersData = await ordersRes.json();

                setOrderTypes(typesData || []);
                setOrders(ordersData || []);
                console.log(typesData );
            } catch (error) {
                console.error('Error fetching data:', error);
            } finally {
                setLoading(false);
            }
        };

        fetchData();
    }, []);

    const handleAddOrderType = async () => {
        if (!newOrderType.trim()) return;

        try {
            const response = await fetch('http://localhost:5000/admin/order-types/add', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                credentials: 'include',
                body: JSON.stringify({ name: newOrderType }),
            });

            const data = await response.json();
            setOrderTypes(prev => [...prev, data]);
            setNewOrderType('');
        } catch (error) {
            console.error('Error adding order type:', error);
        }
    };

    const isValidDate = (dateString) => {
        const regEx = /^\d{4}-\d{2}-\d{2}$/;
        if (!dateString.match(regEx)) return false;
        const d = new Date(dateString);
        return d instanceof Date && !isNaN(d);
    };

    const handleAddOrder = async () => {
        if (!isValidDate(newOrder.date)) {
            setError("Укажите корректную дату в формате ГГГГ-ММ-ДД");
            return;
        }

        if (!newOrder.orderTypeId) {
            setError("Выберите тип приказа");
            return;
        }

        try {
            const response = await fetch('http://localhost:5000/admin/orders/add', {
                method: 'POST',
                mode: 'cors',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json'
                },
                body: JSON.stringify({
                    number: newOrder.number,
                    date: newOrder.date,
                    orderTypeId: newOrder.orderTypeId,
                    description: newOrder.description
                }),
            });

            const data = await response.json();

            if (!response.ok) {
                throw new Error(data.error || 'Ошибка сервера');
            }

            setOrders(prev => [...prev, data]);
            setNewOrder({
                number: '',
                date: new Date().toISOString().slice(0, 10),
                orderTypeId: '',
                description: ''
            });
            setError('');
        } catch (error) {
            console.error('Error adding order:', error);
            setError(error.message);
        }
    };
    if (loading) return <div>Loading...</div>;

    return (
        <div className="admin-orders-container">
            {error && <div className="error-message">{error}</div>}
            <h2>Управление приказами</h2>

            <div className="order-types-section">
                <h3>Типы приказов</h3>
                <div className="add-order-type">
                    <input
                        type="text"
                        value={newOrderType}
                        onChange={(e) => setNewOrderType(e.target.value)}
                        placeholder="Новый тип приказа"
                    />
                    <button onClick={handleAddOrderType}>Добавить</button>
                </div>
                <ul>
                    {orderTypes.map(type => (
                        <li key={type._id}>{type.name}</li>
                    ))}
                </ul>
            </div>

            <div className="orders-section">
                <h3>Создание приказа</h3>
                <div className="order-form">
                    <input
                        type="text"
                        value={newOrder.number}
                        onChange={(e) => setNewOrder({ ...newOrder, number: e.target.value })}
                        placeholder="Номер приказа"
                    />
                    <input
                        type="date"
                        value={newOrder.date}
                        onChange={(e) => setNewOrder({ ...newOrder, date: e.target.value })}
                    />
                    <select
                        value={newOrder.orderTypeId}
                        onChange={(e) => setNewOrder({ ...newOrder, orderTypeId: e.target.value })}
                    >
                        <option value="">Выберите тип приказа</option>
                        {orderTypes.map(type => (
                            <option key={type.id} value={type.id}>{type.name}</option>
                        ))}
                    </select>
                    <textarea
                        value={newOrder.description}
                        onChange={(e) => setNewOrder({ ...newOrder, description: e.target.value })}
                        placeholder="Описание приказа"
                    />
                    <button onClick={handleAddOrder}>Создать приказ</button>
                </div>

                <h3>Список приказов</h3>
                <table>
                    <thead>
                        <tr>
                            <th>Номер</th>
                            <th>Дата</th>
                            <th>Тип</th>
                            <th>Описание</th>
                        </tr>
                    </thead>
                    <tbody>
                        {orders.map(order => (
                            <tr key={order._id}>
                                <td>{order.number}</td>
                                <td>{new Date(order.date).toLocaleDateString()}</td>
                                <td>{order.orderType}</td>
                                <td>{order.description}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default AdminOrders;