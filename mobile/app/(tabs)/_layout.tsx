import { Tabs } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import { Colors } from '../../src/constants/theme';

export default function TabsLayout() {
    return (
        <Tabs
            screenOptions={{
                headerShown: false,
                tabBarActiveTintColor: '#ffffff',
                tabBarInactiveTintColor: '#666666',
                tabBarStyle: {
                    backgroundColor: Colors.backgroundSecondary,
                    borderTopColor: Colors.border,
                    paddingTop: 10,
                    paddingBottom: 10,
                    height: 70,
                },
                tabBarShowLabel: false,
            }}
        >
            <Tabs.Screen
                name="home"
                options={{
                    tabBarIcon: ({ color, size }) => <Ionicons name="home" size={size} color={color} />,
                }}
            />
            <Tabs.Screen
                name="chat"
                options={{
                    href: null, // Remove da navegação
                }}
            />
            <Tabs.Screen
                name="profile"
                options={{
                    tabBarIcon: ({ color, size }) => <Ionicons name="person" size={size} color={color} />,
                }}
            />
            <Tabs.Screen
                name="words"
                options={{
                    href: null, // Remove da navegação
                }}
            />
        </Tabs>
    );
}
