// サイドメニュー機能

// ユーザーの権限に基づいてメニューアイテムを生成
async function createSideMenu() {
    try {
        // ユーザー権限情報を取得
        const response = await fetch('/api/user-permissions');
        const permissions = await response.json();
        
        const menuItems = [
            { id: 'home', name: 'メインメニュー', url: '/', icon: '🏠', screenId: null },
            { id: 'budget', name: '予算管理', url: '/budget', icon: '💰', screenId: '001' },
            { id: 'sales', name: '販売管理', url: '/sales', icon: '📊', screenId: '002' },
            { id: 'inventory', name: '在庫管理', url: '/inventory', icon: '📦', screenId: '003' },
            { id: 'admin', name: '権限管理', url: '/admin', icon: '🔐', screenId: '004' },
            { id: 'user', name: 'ユーザー管理', url: '/user', icon: '👥', screenId: '005' }
        ];
        
        // 現在のパスを取得
        const currentPath = window.location.pathname;
        
        let menuHtml = `
            <div class="side-menu">
                <div class="side-menu-header">
                    <h3>メニュー</h3>
                    <button class="menu-toggle" onclick="toggleSideMenu()">☰</button>
                </div>
                <nav class="side-menu-nav">
                    <ul>
        `;
        
        menuItems.forEach(item => {
            // 権限チェック（メインメニューは常に表示）
            if (item.screenId === null || permissions.screens[item.screenId] === true) {
                const isActive = currentPath === item.url ? 'active' : '';
                menuHtml += `
                    <li>
                        <a href="${item.url}" class="menu-item ${isActive}" data-tooltip="${item.name}">
                            <span class="menu-icon">${item.icon}</span>
                            <span class="menu-text">${item.name}</span>
                        </a>
                    </li>
                `;
            }
        });
        
        menuHtml += `
                    </ul>
                </nav>
            </div>
        `;
        
        return menuHtml;
    } catch (error) {
        console.error('サイドメニューの作成に失敗:', error);
        // エラー時は最小限のメニューを表示
        return `
            <div class="side-menu">
                <div class="side-menu-header">
                    <h3>メニュー</h3>
                    <button class="menu-toggle" onclick="toggleSideMenu()">☰</button>
                </div>
                <nav class="side-menu-nav">
                    <ul>
                        <li><a href="/" class="menu-item" data-tooltip="メインメニュー"><span class="menu-icon">🏠</span><span class="menu-text">メインメニュー</span></a></li>
                    </ul>
                </nav>
            </div>
        `;
    }
}

// サイドメニューの開閉
function toggleSideMenu() {
    const sideMenu = document.querySelector('.side-menu');
    const mainContent = document.querySelector('.main-content');
    
    if (sideMenu) {
        sideMenu.classList.toggle('collapsed');
        if (mainContent) {
            mainContent.classList.toggle('expanded');
        }
        
        // 状態を保存
        const isCollapsed = sideMenu.classList.contains('collapsed');
        localStorage.setItem('sideMenuCollapsed', isCollapsed);
    }
}

// 保存された状態を復元
function restoreSideMenuState() {
    const isCollapsed = localStorage.getItem('sideMenuCollapsed') === 'true';
    const sideMenu = document.querySelector('.side-menu');
    const mainContent = document.querySelector('.main-content');
    
    if (isCollapsed && sideMenu) {
        sideMenu.classList.add('collapsed');
        if (mainContent) {
            mainContent.classList.add('expanded');
        }
    }
}

// ページ読み込み時にサイドメニューを挿入
document.addEventListener('DOMContentLoaded', async function() {
    // body直下に挿入するためのラッパーを作成
    const wrapper = document.createElement('div');
    wrapper.className = 'app-wrapper';
    
    // 既存のコンテンツをラップ
    const bodyContent = Array.from(document.body.children);
    const mainContent = document.createElement('div');
    mainContent.className = 'main-content';
    
    bodyContent.forEach(element => {
        mainContent.appendChild(element);
    });
    
    // サイドメニューを作成
    const sideMenuHtml = await createSideMenu();
    wrapper.innerHTML = sideMenuHtml;
    wrapper.appendChild(mainContent);
    
    // bodyに追加
    document.body.appendChild(wrapper);
    
    // 保存された状態を復元
    restoreSideMenuState();
});

// サイドメニュー用のCSS
const sideMenuStyle = document.createElement('style');
sideMenuStyle.textContent = `
    body {
        margin: 0;
        padding: 0;
    }
    
    .app-wrapper {
        display: flex;
        min-height: 100vh;
        position: relative;
    }
    
    .side-menu {
        width: 250px;
        background-color: #f8f9fa;
        color: #333;
        position: fixed;
        height: 100vh;
        left: 0;
        top: 0;
        transition: width 0.3s ease;
        z-index: 1000;
        overflow-y: auto;
        box-shadow: 2px 0 5px rgba(0,0,0,0.1);
        border-right: 1px solid #e0e0e0;
    }
    
    .side-menu.collapsed {
        width: 60px;
    }
    
    .side-menu-header {
        padding: 20px;
        border-bottom: 1px solid #e0e0e0;
        display: flex;
        justify-content: space-between;
        align-items: center;
        background-color: white;
    }
    
    .side-menu-header h3 {
        margin: 0;
        font-size: 18px;
        color: #2c3e50;
    }
    
    .menu-toggle {
        background: none;
        border: none;
        color: #2c3e50;
        font-size: 24px;
        cursor: pointer;
        padding: 5px;
        transition: transform 0.3s ease;
    }
    
    .menu-toggle:hover {
        transform: scale(1.1);
    }
    
    .side-menu-nav {
        padding: 20px 0;
    }
    
    .side-menu-nav ul {
        list-style: none;
        padding: 0;
        margin: 0;
    }
    
    .menu-item {
        display: flex;
        align-items: center;
        padding: 12px 20px;
        color: #555;
        text-decoration: none;
        transition: all 0.3s ease;
        border-left: 4px solid transparent;
    }
    
    .menu-item:hover {
        background-color: #e8f4f8;
        color: #2c3e50;
        border-left-color: #3498db;
    }
    
    .menu-item.active {
        background-color: #e8f4f8;
        color: #2c3e50;
        border-left: 4px solid #3498db;
        font-weight: 600;
    }
    
    .menu-icon {
        font-size: 20px;
        margin-right: 12px;
        display: inline-block;
        width: 24px;
        text-align: center;
    }
    
    .menu-text {
        font-size: 14px;
        color: inherit;
    }
    
    .main-content {
        margin-left: 250px;
        flex: 1;
        transition: margin-left 0.3s ease;
        min-width: 0;
        width: calc(100% - 250px);
    }
    
    .main-content.expanded {
        margin-left: 60px;
        width: calc(100% - 60px);
    }
    
    .side-menu.collapsed .menu-text {
        display: none;
    }
    
    .side-menu.collapsed .side-menu-header h3 {
        display: none;
    }
    
    .side-menu.collapsed .menu-item {
        justify-content: center;
        padding: 12px 10px;
    }
    
    .side-menu.collapsed .menu-icon {
        margin-right: 0;
        font-size: 24px;
    }
    
    /* ツールチップ */
    .side-menu.collapsed .menu-item {
        position: relative;
    }
    
    .side-menu.collapsed .menu-item:hover::after {
        content: attr(data-tooltip);
        position: absolute;
        left: 100%;
        top: 50%;
        transform: translateY(-50%);
        background-color: #2c3e50;
        color: white;
        padding: 5px 10px;
        border-radius: 4px;
        white-space: nowrap;
        font-size: 14px;
        margin-left: 10px;
        z-index: 1001;
        box-shadow: 0 2px 5px rgba(0,0,0,0.2);
    }
    
    /* コンテナの最大幅を調整 */
    .main-content .container {
        max-width: 100%;
        padding: 20px;
    }
    
    /* ヘッダーの調整 */
    .main-content .header {
        margin: 0;
        position: sticky;
        top: 0;
        z-index: 10;
    }
    
    /* ユーザー切り替えボタンとの干渉を防ぐ */
    #user-switcher {
        z-index: 1100;
    }
    
    /* モバイル対応 */
    @media (max-width: 768px) {
        .side-menu {
            width: 60px;
        }
        
        .side-menu.collapsed {
            width: 250px;
        }
        
        .main-content {
            margin-left: 60px;
            width: calc(100% - 60px);
        }
        
        .main-content.expanded {
            margin-left: 250px;
            width: calc(100% - 250px);
        }
        
        .side-menu:not(.collapsed) .menu-text {
            display: block;
        }
        
        .side-menu:not(.collapsed) .menu-item {
            justify-content: flex-start;
        }
        
        .side-menu:not(.collapsed) .menu-icon {
            margin-right: 12px;
            font-size: 20px;
        }
        
        .menu-toggle {
            display: block;
        }
    }
`;
document.head.appendChild(sideMenuStyle);