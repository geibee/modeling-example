// 共通のユーザー情報表示関数

// ユーザー情報を取得して表示する共通関数
async function displayUserInfo() {
    try {
        const response = await fetch('/api/session');
        if (response.ok) {
            const data = await response.json();
            if (data.authenticated) {
                // ユーザーIDを表示
                const userIdElement = document.getElementById('user-id');
                if (userIdElement) {
                    userIdElement.textContent = data.user_id;
                }
                
                // ロールと組織情報を取得
                try {
                    const rolesResponse = await fetch('/api/current-user');
                    if (rolesResponse.ok) {
                        const userData = await rolesResponse.json();
                        
                        // ロール名を表示
                        const rolesElement = document.getElementById('roles');
                        if (rolesElement) {
                            rolesElement.textContent = userData.role_name || 'Unknown';
                        }
                        
                        // 組織情報を表示
                        const orgElement = document.getElementById('organizations');
                        if (orgElement) {
                            if (userData.organizations && userData.organizations.length > 0) {
                                const orgNames = userData.organizations.map(org => org.org_name).join(', ');
                                orgElement.textContent = orgNames;
                            } else {
                                orgElement.textContent = '所属なし';
                            }
                        }
                    } else {
                        // エラーの場合でもデフォルト値を設定
                        const rolesElement = document.getElementById('roles');
                        if (rolesElement) {
                            rolesElement.textContent = 'Unknown';
                        }
                        const orgElement = document.getElementById('organizations');
                        if (orgElement) {
                            orgElement.textContent = '取得エラー';
                        }
                    }
                } catch (error) {
                    console.error('ユーザー詳細情報の取得エラー:', error);
                    // エラーの場合でもデフォルト値を設定
                    const rolesElement = document.getElementById('roles');
                    if (rolesElement) {
                        rolesElement.textContent = 'Unknown';
                    }
                    const orgElement = document.getElementById('organizations');
                    if (orgElement) {
                        orgElement.textContent = '取得エラー';
                    }
                }
            }
        }
    } catch (error) {
        console.error('ユーザー情報の取得エラー:', error);
        const userIdElement = document.getElementById('user-id');
        if (userIdElement) {
            userIdElement.textContent = 'Unknown';
        }
        const rolesElement = document.getElementById('roles');
        if (rolesElement) {
            rolesElement.textContent = 'Unknown';
        }
        const orgElement = document.getElementById('organizations');
        if (orgElement) {
            orgElement.textContent = '取得エラー';
        }
    }
}