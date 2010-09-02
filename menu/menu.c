#include <gtk/gtk.h>
#include <stdio.h>


/*
 *  メニューの中のどれかが選択されたときに呼ばれるコールバック関数。
 *  label にセットされているテキストを変更する。
 */
static void menu_activated(GtkMenuItem *menuitem, GtkWidget *label)
{
    const gchar *name = gtk_menu_item_get_label(menuitem);
    char message[100];
    snprintf(message, 100, "%sの注文を承りました。", name);

    gtk_label_set_text(GTK_LABEL(label), message); 
}


int main(int argc, char *argv[]) 
{
    GtkWidget *window, *vbox, *label;
    GtkWidget *menubar, *lunch, *lunchmenu;
    GtkWidget *menu1, *menu2, *menu3;

    gtk_init(&argc, &argv);

    window = gtk_window_new(GTK_WINDOW_TOPLEVEL);
    gtk_window_set_title(GTK_WINDOW(window), "Menubar");
    gtk_widget_set_size_request(window, 300, 200);

    g_signal_connect(G_OBJECT(window), "destroy",
                     G_CALLBACK(gtk_main_quit), NULL);

    //  MenuBar と Label を格納するための BOX
    vbox = gtk_vbox_new(FALSE, 5);
    
    //  ラベル。
    //  メニューを選択するとこのテキストが変わる
    label = gtk_label_new("メニューを選んでください");
    
    //  (1)MenuBar 'menubar' を作成する
    menubar = gtk_menu_bar_new();
    
    //  (2)MenuItem 'lunch' を作成する
    lunch = gtk_menu_item_new_with_label("メニュー");
    
    //  (3)Menu 'lunchmenu' を作成する
    lunchmenu = gtk_menu_new();

    //  (4a)MenuItem 'menu*' を作成する
    menu1 = gtk_menu_item_new_with_label("カレー");
    menu2 = gtk_menu_item_new_with_label("おにぎり");
    menu3 = gtk_menu_item_new_with_label("俺の塩");

    //  (4b)MenuItem 'menu*' を Menu 'lunchmenu' に追加する
    gtk_menu_shell_append(GTK_MENU_SHELL(lunchmenu), menu1);
    gtk_menu_shell_append(GTK_MENU_SHELL(lunchmenu), menu2);
    gtk_menu_shell_append(GTK_MENU_SHELL(lunchmenu), menu3);

    //  (5)MenuItem 'menu*' が選択されたときに呼ばれるコールバック関数を
    //     menu_activated() に設定する。
    //     また、このときの引数として label も渡すようにする
    g_signal_connect(G_OBJECT(menu1), "activate",
                     G_CALLBACK(menu_activated), (gpointer)label);
    g_signal_connect(G_OBJECT(menu2), "activate", 
                     G_CALLBACK(menu_activated), (gpointer)label);
    g_signal_connect(G_OBJECT(menu3), "activate",
                     G_CALLBACK(menu_activated), (gpointer)label);
    
    //  (3)Menu 'lunchmenu' を MenuItem 'lunch' のサブメニューとする
    gtk_menu_item_set_submenu(GTK_MENU_ITEM(lunch), lunchmenu);
    
    //  (2)MenuItem 'lunch' を MenuBar 'menubar' に追加する
    gtk_menu_shell_append(GTK_MENU_SHELL(menubar), lunch);

    //  (1)MenuBar 'menubar' を BOX に追加する
    gtk_box_pack_start(GTK_BOX(vbox), menubar, FALSE, FALSE, 0);

    gtk_box_pack_start(GTK_BOX(vbox), label, TRUE, TRUE, 0);
    gtk_container_add(GTK_CONTAINER(window), vbox);
    gtk_widget_show_all(window);

    gtk_main();

    return 0;
}


