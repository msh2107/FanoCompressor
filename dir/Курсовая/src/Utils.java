import java.io.File;
import java.nio.file.Files;
import java.security.MessageDigest;
import java.security.SecureRandom;
import java.util.ArrayList;
import java.util.HexFormat;

public class Utils {

    private static String password;
    private static File initVectorFile;
    private static ArrayList<File> files;

    public static boolean processInput(String[] args) throws Exception {
        if (!(args != null && (args.length == 1 && "test".equals(args[0]) || args.length >= 3))) {
            throw new Exception("Wrong arguments!");
        }

        if ("test".equals(args[0])) {
            return true;
        } else {
            password = args[0];
            initVectorFile = new File(args[1]);
            files = new ArrayList<>();

            for (int i = 2; i < args.length; i++) {
                File file = new File(args[i]);
                if (!file.isDirectory()) {
                    files.add(file);
                } else {
                    addFilesFromDirectory(file, files);
                }
            }
        }

        return false;
    }

    private static void addFilesFromDirectory(File directory, ArrayList<File> files) {
        File[] directoryFiles = directory.listFiles();

        if (directoryFiles != null) {
            for (File directoryFile : directoryFiles) {
                if (!directoryFile.isDirectory()) {
                    files.add(directoryFile);
                } else {
                    addFilesFromDirectory(directoryFile, files);
                }
            }
        }
    }

    public static String passwordToKey() throws Exception {
        String key = HexFormat.of().formatHex(
                MessageDigest.getInstance("SHA-256").digest(password.getBytes()));

        return key.substring(48);
    }

    public static String getInitVector() throws Exception {
        if (initVectorFile.length() == 8L)
            return HexFormat.of().formatHex(
                    Files.readAllBytes(initVectorFile.toPath()));
        else {
            byte[] initVector = new byte[8];

            SecureRandom secureRandom = new SecureRandom();
            secureRandom.nextBytes(initVector);

            Files.write(initVectorFile.toPath(), initVector);

            return HexFormat.of().formatHex(initVector);
        }
    }

    public static ArrayList<File> getFiles() {
        return files;
    }
}
