// ユーザー切り替え機能

let users = [
    { id: 'user001', name: '管理職ユーザー', role: '管理職 (R01)', permissions: '予算管理' },
    { id: 'user002', name: '営業担当ユーザー', role: '営業担当 (R02)', permissions: '販売管理' },
    { id: 'user003', name: '在庫管理者ユーザー', role: '在庫管理者 (R03)', permissions: '在庫管理' },
    { id: 'admin', name: 'システム管理者', role: 'システム管理者 (R04)', permissions: '全権限' }
];

// ユーザー一覧を動的に取得
async function loadUserList() {
    try {
        const response = await fetch('/api/users-list');
        if (response.ok) {
            const userList = await response.json();
            // ユーザーリストを更新
            users = userList.map(user => ({
                id: user.user_id,
                name: getUserDisplayName(user.user_id),
                role: `${user.role_name || getRoleName(user.role_id)} (${user.role_id})`,
                permissions: getRolePermissions(user.role_id)
            }));
        }
    } catch (error) {
        console.log('ユーザー一覧の取得に失敗、デフォルトのリストを使用');
        // エラーでもデフォルトのユーザーリストは既に定義されているので使用可能
    }
}

// getUserDisplayNameを外部で定義（上部に移動）
function getUserDisplayName(userID) {
    switch (userID) {
        case 'user001':
            return '管理職ユーザー';
        case 'user002':
            return '営業担当ユーザー';
        case 'user003':
            return '在庫管理者ユーザー';
        case 'admin':
            return 'システム管理者';
        default:
            return userID; // フォールバックとしてユーザーIDを表示
    }
}

// ロールに基づく権限の説明を取得
function getRolePermissions(roleId) {
    const permissions = {
        'R01': '予算管理',
        'R02': '販売管理',
        'R03': '在庫管理',
        'R04': '全権限'
    };
    return permissions[roleId] || '不明';
}

// ロールIDからロール名を取得
function getRoleName(roleId) {
    const roleNames = {
        'R01': '管理職',
        'R02': '営業担当',
        'R03': '在庫管理者',
        'R04': 'システム管理者'
    };
    return roleNames[roleId] || '不明';
}

// 現在のセッション情報を取得
async function getCurrentSessionInfo() {
    try {
        const response = await fetch('/api/session', {
            headers: {
                'Accept': 'application/json'
            }
        });
        
        // レスポンスのタイプを確認
        const contentType = response.headers.get('content-type');
        if (!contentType || !contentType.includes('application/json')) {
            console.error('セッションAPIが正しいJSONを返していません');
            return { authenticated: false };
        }
        
        if (response.ok) {
            const data = await response.json();
            if (data.authenticated) {
                return {
                    authenticated: true,
                    user_id: data.user_id,
                    roles: data.roles
                };
            }
        }
    } catch (error) {
        console.error('セッション情報の取得に失敗:', error);
    }
    return { authenticated: false };
}

// 現在のユーザー情報をAPIから取得
async function getCurrentUserInfo() {
    const session = await getCurrentSessionInfo();
    if (!session.authenticated) {
        return { id: 'guest', name: 'ゲスト', role: '未認証', permissions: 'なし' };
    }
    
    try {
        const response = await fetch('/api/current-user', {
            headers: {
                'Accept': 'application/json'
            }
        });
        
        // レスポンスのタイプを確認
        const contentType = response.headers.get('content-type');
        if (contentType && contentType.includes('application/json') && response.ok) {
            const data = await response.json();
            return {
                id: data.user_id,
                name: data.name,
                role: `${data.role_name} (${data.roles[0]})`,
                permissions: data.permissions
            };
        }
    } catch (error) {
        console.error('ユーザー情報の取得に失敗:', error);
    }
    
    // フォールバック
    const user = users.find(u => u.id === session.user_id);
    return user || { id: session.user_id, name: session.user_id, role: '不明', permissions: '不明' };
}

// ユーザー切り替えUIを作成
async function createUserSwitcher() {
    const currentUser = await getCurrentUserInfo();
    const currentUserId = currentUser.id;
    
    const switcherHtml = `
        <div class="user-switcher">
            <div class="current-user" onclick="toggleUserDropdown()">
                <div class="user-info">
                    <span class="user-name">${currentUser.name}</span>
                    <span class="user-role">${currentUser.role}</span>
                </div>
                <svg class="dropdown-arrow" width="12" height="12" viewBox="0 0 12 12">
                    <path d="M2 4L6 8L10 4" stroke="currentColor" stroke-width="2" fill="none"/>
                </svg>
            </div>
            <div class="user-dropdown" id="userDropdown">
                <div class="dropdown-header">ユーザー切り替え</div>
                ${users.map(user => `
                    <div class="user-option ${user.id === currentUserId ? 'active' : ''}" onclick="switchUser('${user.id}')">
                        <div class="user-option-info">
                            <div class="user-option-name">${user.name}</div>
                            <div class="user-option-role">${user.role}</div>
                            <div class="user-option-permissions">権限: ${user.permissions}</div>
                        </div>
                        ${user.id === currentUserId ? '<span class="check-mark">✓</span>' : ''}
                    </div>
                `).join('')}
            </div>
        </div>
    `;
    
    return switcherHtml;
}

// ドロップダウンの表示/非表示を切り替え
function toggleUserDropdown() {
    const dropdown = document.getElementById('userDropdown');
    dropdown.classList.toggle('show');
    
    // クリック外で閉じる
    document.addEventListener('click', function closeDropdown(e) {
        if (!e.target.closest('.user-switcher')) {
            dropdown.classList.remove('show');
            document.removeEventListener('click', closeDropdown);
        }
    });
}

// ユーザーを切り替え
async function switchUser(userId) {
    try {
        // まずログアウト
        await fetch('/api/logout', { method: 'POST' });
        
        // 新しいユーザーでログイン
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ user_id: userId })
        });
        
        if (response.ok) {
            // ログイン成功 - 現在のページをリロード
            window.location.reload();
        } else {
            alert('ユーザー切り替えに失敗しました');
        }
    } catch (error) {
        console.error('ユーザー切り替えエラー:', error);
        alert('ユーザー切り替え中にエラーが発生しました');
    }
}

// ユーザー切り替えUIを初期化する関数
async function initializeUserSwitcher() {
    console.log('ユーザー切り替えUIを初期化中...');
    
    // 最初にユーザーリストを取得
    await loadUserList();
    
    // ヘッダーが存在する場合、その中に挿入
    const header = document.querySelector('.header');
    if (header) {
        console.log('ヘッダー要素が見つかりました');
        
        // 既存のユーザー切り替えUIがある場合は削除
        const existingSwitcher = header.querySelector('.user-switcher');
        if (existingSwitcher) {
            existingSwitcher.parentElement.remove();
        }
        
        const switcherContainer = document.createElement('div');
        switcherContainer.innerHTML = await createUserSwitcher();
        header.appendChild(switcherContainer);
        console.log('ユーザー切り替えUIを追加しました');
    } else {
        console.error('ヘッダー要素が見つかりません。ページ構造を確認してください。');
        // 少し待ってから再試行
        setTimeout(async () => {
            const retryHeader = document.querySelector('.header');
            if (retryHeader) {
                console.log('再試行: ヘッダー要素が見つかりました');
                const switcherContainer = document.createElement('div');
                switcherContainer.innerHTML = await createUserSwitcher();
                retryHeader.appendChild(switcherContainer);
            }
        }, 500);
    }
}

// ページ読み込み時にユーザー切り替えUIを挿入
// DOMContentLoadedがすでに発火している場合も考慮
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initializeUserSwitcher);
} else {
    // DOMはすでに読み込まれている
    initializeUserSwitcher();
}

// CSSを動的に追加
const style = document.createElement('style');
style.textContent = `
    .header {
        position: relative;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }
    
    .user-switcher {
        position: absolute;
        right: 20px;
        top: 50%;
        transform: translateY(-50%);
    }
    
    .current-user {
        background: rgba(255, 255, 255, 0.1);
        border: 1px solid rgba(255, 255, 255, 0.2);
        border-radius: 6px;
        padding: 8px 16px;
        cursor: pointer;
        display: flex;
        align-items: center;
        gap: 12px;
        transition: all 0.3s ease;
        min-width: 200px;
    }
    
    .current-user:hover {
        background: rgba(255, 255, 255, 0.2);
        border-color: rgba(255, 255, 255, 0.3);
    }
    
    .user-info {
        flex: 1;
        text-align: left;
    }
    
    .user-name {
        display: block;
        color: white;
        font-weight: 500;
        font-size: 14px;
    }
    
    .user-role {
        display: block;
        color: rgba(255, 255, 255, 0.8);
        font-size: 12px;
        margin-top: 2px;
    }
    
    .dropdown-arrow {
        transition: transform 0.3s ease;
        color: white;
    }
    
    .user-dropdown {
        position: absolute;
        top: 100%;
        right: 0;
        margin-top: 8px;
        background: white;
        border-radius: 8px;
        box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
        min-width: 280px;
        opacity: 0;
        visibility: hidden;
        transform: translateY(-10px);
        transition: all 0.3s ease;
        z-index: 10000;
    }
    
    .user-dropdown.show {
        opacity: 1;
        visibility: visible;
        transform: translateY(0);
    }
    
    .user-dropdown.show ~ .current-user .dropdown-arrow {
        transform: rotate(180deg);
    }
    
    .dropdown-header {
        padding: 12px 16px;
        font-weight: 600;
        color: #333;
        border-bottom: 1px solid #e0e0e0;
        font-size: 14px;
    }
    
    .user-option {
        padding: 12px 16px;
        cursor: pointer;
        transition: background-color 0.2s ease;
        display: flex;
        align-items: center;
        justify-content: space-between;
    }
    
    .user-option:hover {
        background-color: #f5f7fa;
    }
    
    .user-option.active {
        background-color: #e8f4f8;
    }
    
    .user-option-info {
        flex: 1;
    }
    
    .user-option-name {
        font-weight: 500;
        color: #333;
        font-size: 14px;
    }
    
    .user-option-role {
        color: #666;
        font-size: 12px;
        margin-top: 2px;
    }
    
    .user-option-permissions {
        color: #999;
        font-size: 11px;
        margin-top: 2px;
    }
    
    .check-mark {
        color: #3498db;
        font-weight: bold;
        font-size: 18px;
    }
    
    @media (max-width: 768px) {
        .user-switcher {
            position: static;
            transform: none;
            margin-top: 10px;
            width: 100%;
        }
        
        .current-user {
            width: 100%;
        }
        
        .user-dropdown {
            left: 0;
            right: 0;
            width: auto;
        }
        
        .header {
            flex-direction: column;
            padding-bottom: 60px;
        }
    }
`;
document.head.appendChild(style);