#include <stdio.h>
#include <stdlib.h>
#include <libxml/xmlreader.h>


//  現在パース中のノード名
typedef enum {
    STATE_NONE,
    STATE_ITEM,
    STATE_PRICE
} parsingStatus;


//  1つのノードを処理する
void processNode(xmlTextReaderPtr reader) 
{
    static parsingStatus state = STATE_NONE;
    xmlElementType nodeType;
    xmlChar *name, *value;

    //  ノード情報の取得
    nodeType = xmlTextReaderNodeType(reader);       //  ノードタイプ
    name = xmlTextReaderName(reader);               //  ノード名
    if (!name) 
        name = xmlStrdup(BAD_CAST "---");

    if (nodeType == XML_READER_TYPE_ELEMENT) {              //  開始
        if ( xmlStrcmp(name, BAD_CAST "item") == 0 ) {
            state = STATE_ITEM;

        } else if ( xmlStrcmp(name, BAD_CAST "price") == 0 ) {
            state = STATE_PRICE;
        }

    } else if (nodeType == XML_READER_TYPE_END_ELEMENT) {   //  終了
        if ( xmlStrcmp(name, BAD_CAST "fruit") == 0 ) {
            printf("-----------------------\n"); 
        }
        
        state = STATE_NONE;

    } else if (nodeType == XML_READER_TYPE_TEXT) {          //  テキスト
        //  テキストを取得する
        value = xmlTextReaderValue(reader);
        
        if (!value)
            value = xmlStrdup(BAD_CAST "---");

        if ( state == STATE_ITEM ) {
            printf("品名: %s\n", value);

        } else if ( state == STATE_PRICE ) {
            printf("価格: %s円\n", value);
        }

        xmlFree(value);
    }

    xmlFree(name);
}



int main(int argc, char *argv[])
{
    xmlTextReaderPtr reader;
    int ret;

    //  Readerの作成
    reader = xmlNewTextReaderFilename("./sample.xml");
    if ( !reader ) {
        printf("Failed to open XML file.\n");
        return 1;
    }

    printf("-----------------------\n"); 
    
    //  次のノードに移動 
    ret = xmlTextReaderRead(reader);
    while (ret == 1) {
        //  現在のノードを処理
        processNode(reader);

        //  次のノードに移動
        ret = xmlTextReaderRead(reader);
    }
    
    //  Reader のすべてのリソースを開放
    xmlFreeTextReader(reader);

    //  xmlTextReaderRead の戻り値が -1 だった場合はパースエラー
    if (ret == -1) {
        printf("Parse error.\n");
        return 1;
    }

    return 0;
}





