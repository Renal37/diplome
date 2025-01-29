const AdminCheckProfile = () => {
    const [users, setUsers] = useState([]);
    
    useEffect(() => {
        fetch("http://localhost:5000/users")
            .then(res => res.json())
            .then(data => setUsers(data));
    }, []);
    
    return (
        <div>
            <h1>Проверка профилей</h1>
            {users.map(user => (
                <div key={user._id}>
                    <p>{user.fullName} ({user.username})</p>
                    <img src={user.passportData} alt="Паспорт" width={100} />
                    <img src={user.snils} alt="СНИЛС" width={100} />
                    <button>Подтвердить</button>
                </div>
            ))}
        </div>
    );
};

export default AdminCheckProfile