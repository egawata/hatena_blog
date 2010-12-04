#include <gtk/gtk.h>


/* 
 *  Model の定義
 */
void create_model(GtkWidget *treeview)
{
    GtkListStore *store;

    /*  GtkListStore の生成、および列定義  */
    store = gtk_list_store_new(
                4,              //  要素の個数
                G_TYPE_STRING,  //  場所
                G_TYPE_STRING,  //  天気
                G_TYPE_INT,     //  最高気温
                G_TYPE_INT      //  最低気温
            );

    /*  treeview に結びつける  */
    gtk_tree_view_set_model(GTK_TREE_VIEW(treeview), GTK_TREE_MODEL(store)); 

    /*  一旦結びつけたら store 自体は開放してよい */
    g_object_unref(store);
}


/*
 * Model に実データ追加
 */
void set_data(GtkWidget *treeview)
{
    GtkListStore *store;
    GtkTreeIter iter;

    //  treeview に結びつけられている Model を取得する。
    //  一旦、中のデータをすべて消去する。
    store = GTK_LIST_STORE( gtk_tree_view_get_model( GTK_TREE_VIEW(treeview) ) );
    gtk_list_store_clear(store);

    /*  新しいレコードを生成する(3レコード) */
    gtk_list_store_append(store, &iter);
    gtk_list_store_set(store, &iter, 0, "東京", 1, "晴のち曇", 2, 17, 3, 9, -1);

    gtk_list_store_append(store, &iter);
    gtk_list_store_set(store, &iter, 0, "群馬", 1, "晴時々曇", 2, 18, 3, 5, -1);

    gtk_list_store_append(store, &iter);
    gtk_list_store_set(store, &iter, 0, "大阪", 1, "晴れ", 2, 16, 3, 7, -1);
}



/*
 * treeview に列を1つ追加
 */
void append_column_to_treeview(GtkWidget *treeview, const char *title, const int order)
{
    GtkTreeViewColumn   *column;
    GtkCellRenderer *renderer;

    /*   CellRenderer の作成  */
    renderer = gtk_cell_renderer_text_new();
    
    /*
     *   列を作成し、CellRenderer を追加する。 
     *   さらに Model との関連を定義する
     */
    column = gtk_tree_view_column_new_with_attributes(
                title,          //  ヘッダ行に表示するタイトル
                renderer,       //  CellRenderer
                "text",  order, //  Model の order 番めのデータを、renderer の属性 "text" に結びつける
                NULL            //  終端
             );
    
    /*   今生成した列を、大元の treeview に追加する  */
    gtk_tree_view_append_column(GTK_TREE_VIEW(treeview), column);
}


/*
 * treeview に View を定義する
 */
void create_view(GtkWidget *treeview)
{
    append_column_to_treeview(treeview, "場所", 0);
    append_column_to_treeview(treeview, "天気", 1);
    append_column_to_treeview(treeview, "最高気温", 2);
    append_column_to_treeview(treeview, "最低気温", 3);    
}
    

int main(int argc, char **argv)
{
    GtkWidget *window, *treeview;
    
    gtk_init(&argc, &argv);

    window = gtk_window_new(GTK_WINDOW_TOPLEVEL);
    gtk_window_set_title(GTK_WINDOW(window), "明日の天気");
    g_signal_connect(G_OBJECT(window), "destroy",
                     G_CALLBACK(gtk_main_quit), NULL);

    /*  GtkTreeView の生成  */
    treeview = gtk_tree_view_new();

    /*  Model の定義  */
    create_model(treeview);

    /*  Model に実データ追加  */
    set_data(treeview);

    /*  View の定義  */
    create_view(treeview);

    gtk_container_add(GTK_CONTAINER(window), treeview);
    gtk_widget_show_all(window);

    gtk_main();
    return 0;
}



