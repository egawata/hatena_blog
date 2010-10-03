#include <gtk/gtk.h>
#include <glib-object.h>
#include <stdio.h>


/*
 *  カレンダー上の日付がクリックされたときに呼び出される
 *  コールバック関数
 */
static void disp_day(GtkCalendar *calendar, GtkEntry *text)
{
    guint year, month, day;
    gchar ymd[100];

    //  カレンダー上で選択された年月日を得る
    gtk_calendar_get_date(calendar, &year, &month, &day);
    month++;    //  月は 0-11 の値で取得されるので +1

    snprintf(ymd, 100, "%d年 %d月 %d日", year, month, day);
    
    gtk_entry_set_text(text, ymd);
}


int main(int argc, char *argv[]) 
{
    GtkWidget *window, *vbox, *text, *calendar;

    gtk_init(&argc, &argv);

    window = gtk_window_new(GTK_WINDOW_TOPLEVEL);
    gtk_window_set_title(GTK_WINDOW(window), "カレンダー");
    
    g_signal_connect(G_OBJECT(window), "destroy",
                     G_CALLBACK(gtk_main_quit), NULL);

    vbox = gtk_vbox_new(FALSE, 5);

    //  GtkCalendar オブジェクトの生成
    calendar = gtk_calendar_new();
    
    //  日付表示用エリア
    text = gtk_entry_new();
    
    //  カレンダー上の日付がクリックされたら
    //  disp_day 関数をコールする
    g_signal_connect(G_OBJECT(calendar), "day-selected", 
                     G_CALLBACK(disp_day), (gpointer)text);

    gtk_box_pack_start_defaults(GTK_BOX(vbox), calendar);
    gtk_box_pack_start_defaults(GTK_BOX(vbox), text);
     
    gtk_container_add(GTK_CONTAINER(window), vbox);

    //  起動直後は選択された日付が text に表示されていないので、
    //  明示的に day-selected シグナルを発生させる
    g_signal_emit_by_name(calendar, "day-selected", text);

    gtk_widget_show_all(window);
    gtk_main();

    return 0;
}


