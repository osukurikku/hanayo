<?php
//Global Variables
$GLOBAL['stream'] = $_GET['releasestream'];
$GLOBAL['key'] = $_GET['key'];
$GLOBAL['conndb'] = conndb();

//Database Settings
$sql_settings['ServerName'] = "";
$sql_settings['User'] = "";
$sql_settings['Password'] = "";
$sql_settings['Database'] = "";

//Key and Stream
$key = $GLOBAL['key'];
$stream = $GLOBAL['stream'];

if($key != "SECRET APIKEY HERE"){
    die("Wrong Key!");
}
if($stream = null){
    die("Wrong releasestream!");
} else {
    switch($stream){
        case "cuttingedge":
            $stream = "cuttingedge";
            break;
        case "stable":
            $stream = "stable";
            break;
        default: 
            $stream = null;
            break;
    }
}
$target_dir = scandir('./dl/');
foreach($target_dir as $file =>$value){ 
    if (!in_array($value,array(".","..","index.php","push.php"))) { 
        $file = pathinfo($value);
        // File Name
        $fn = $file['basename'];
        // File Hash
        $fh = md5_file($fn);
        // File extension
        $ext = $file['extension'];
        // File Size
        if($ext != ".zip"){
            $fs = filesize($fn);
        }
        // Time Stamp
        $ts = date('Y-m-d h:m:s');
        // Out Dir
        $out = "./$stream/".$value;
        // Out ZIP
        $zipfile = $out.".zip";
        // File Zip Size
        $fsz = filesize($out.".zip");

        if(file_exists($out)){
            unlink($out);
        }

        if(updatedb($fn, $fh, $fs, $fsz, $ts)){
            if(rename($fn, $out)){
                echo "File moved!";
            }
            $error = 0;
        } else {
            $error = 1;
        }
    }
}
//Check if ERROR (make no sense)
if($error){
    die("The side died! Sry :(");
} else {
    echo 1;
}

//Connect to database
function conndb()
{
    $servername = $sql_settings['ServerName'];
    $username = $sql_settings['User'];
    $password = $sql_settings['Password'];
    $dbname = $sql_settings['Database'];
    $conn = new mysqli($servername, $username, $password, $dbname);

    if ($sql->connect_error) {
        die("Whoops, looks like someone droped the Server and killed the MySQL Connection. " . $conn->connect_error);
    } 

    return $conn;
}

// Update DATABASE
function updatedb($fn, $fh, $fs, $fsz, $ts)
{
    $conn = conndb();
    $ivs = increaseversion($fn);

    $sql = "UPDATE ".$stream."_check SET file_version='".$ivs."', file_hash='".$fh."', filesize='".$fs."', filesize_zip='".$fsz."', timestamp='".$ts."' WHERE filename='$fn'";
    if ($conn->query($sql) === TRUE) {
        return true;
    } else {
        return false;
    }
}

// Get Current Version
function getversion($file)
{
    $conn = conndb();

    $sql = "SELECT * FROM ".$stream."_table";
    
    $result = $conn->query($sql);
        if ($result->num_rows > 0) {
            while($row = $result->fetch_assoc()){
                $row = auto_cast($row);
                if($row['filename'] == $file){
                    $out = $row['file_version'];
                }
            }
        }
    return $out;
}

// Increase File Version
function increaseversion($file)
{
    $conn = conndb();

    $sql = "SELECT * FROM ".$stream."_table";
    
    $result = $conn->query($sql);
        if ($result->num_rows > 0) {
            while($row = $result->fetch_assoc()){
                if($row['filename'] == $file){
                    $out = ++$row['file_version'];
                }
            }
        }
    return $out;
}